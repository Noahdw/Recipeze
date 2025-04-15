package parsing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractedRecipe represents the extracted recipe data
type ExtractedRecipe struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURL     string   `json:"image_url"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	PrepTime     string   `json:"prep_time,omitempty"`
	CookTime     string   `json:"cook_time,omitempty"`
	TotalTime    string   `json:"total_time,omitempty"`
	Yield        string   `json:"yield,omitempty"`
	Author       string   `json:"author,omitempty"`
}

// JSONLDRecipe represents the JSON-LD schema.org/Recipe structure
type JSONLDRecipe struct {
	Context      string      `json:"@context"`
	Type         string      `json:"@type"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Image        interface{} `json:"image"` // Can be string or array of strings
	Author       interface{} `json:"author"`
	PrepTime     string      `json:"prepTime"`
	CookTime     string      `json:"cookTime"`
	TotalTime    string      `json:"totalTime"`
	RecipeYield  interface{} `json:"recipeYield"`
	Ingredients  interface{} `json:"recipeIngredient"`   // Can be array of strings
	Instructions interface{} `json:"recipeInstructions"` // Can be array of strings or objects
}

// ParseRecipe extracts recipe information from HTML content
func ParseRecipe(htmlContent []byte) (*ExtractedRecipe, error) {
	recipe := &ExtractedRecipe{
		Ingredients:  []string{},
		Instructions: []string{},
	}

	// Try to parse structured JSON-LD data first
	if jsonLDRecipe := extractJSONLD(htmlContent); jsonLDRecipe != nil {
		mapJSONLDToRecipe(jsonLDRecipe, recipe)
		slog.Info("found jsonLD in recipe")
	}

	// Parse HTML for additional information and as fallback
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})

	// Try to extract schema.org microdata if JSON-LD wasn't found
	if len(recipe.Ingredients) == 0 || len(recipe.Instructions) == 0 {
		extractMicrodata(doc, recipe)
	}
	if len(recipe.Ingredients) == 0 {
		return nil, fmt.Errorf("failed to find ingredients")
	}
	return recipe, nil
}

func HtmlToText(htmlContent []byte) []byte {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	// First, remove clearly non-content elements
	nonContentSelectors := []string{
		"script", "style", "noscript", "iframe", "svg", "canvas",
		"header", "footer", "nav", "aside",
		".comments", ".comment-section", "#comments",
		".advertisement", ".ads", ".ad-container",
		".sidebar", ".navigation", ".menu", ".social-media",
		".related-posts", ".recommended", ".cookie-notice",
		".popup", ".modal", ".newsletter-signup",
	}

	for _, selector := range nonContentSelectors {
		doc.Find(selector).Each(func(i int, el *goquery.Selection) {
			el.Remove()
		})
	}

	// Try to identify the main content area if available
	contentSelectors := []string{
		"article", ".post-content", ".entry-content", ".content",
		"#content", ".article-content", ".main-content",
		"[role='main']", ".post-body", ".article-body",
		".wprm-recipe-ingredients-container", ".wprm-recipe-instructions-container",
	}

	var mainContent string
	for _, selector := range contentSelectors {
		selected := doc.Find(selector)
		if selected.Length() > 0 {
			// Found likely content area - use only this section
			mainContent = selected.Text()
			break
		}
	}

	var text string
	if mainContent != "" {
		text = mainContent
	} else {
		// Fallback to the whole document if no main content area identified
		text = doc.Text()
	}

	// Process the text to clean it up
	text = html.UnescapeString(text) // Fix HTML entities

	// Remove excessive whitespace but preserve paragraph structure
	re := regexp.MustCompile(`[ \t]+`)
	text = re.ReplaceAllString(text, " ")

	// Replace multiple newlines with a single newline
	re = regexp.MustCompile(`\n+`)
	text = re.ReplaceAllString(text, "\n")

	// Remove lines with only whitespace
	re = regexp.MustCompile(`(?m)^\s*$`)
	text = re.ReplaceAllString(text, "")

	// Trim trailing whitespace on each line
	re = regexp.MustCompile(`(?m) +$`)
	text = re.ReplaceAllString(text, "")

	// Final trim of leading/trailing whitespace
	text = strings.TrimSpace(text)

	return []byte(text)
}

// extractJSONLD extracts JSON-LD data from HTML
func extractJSONLD(htmlContent []byte) *JSONLDRecipe {
	re := regexp.MustCompile(`<script type="application/ld\+json">(.*?)</script>`)
	matches := re.FindAllSubmatch(htmlContent, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			jsonContent := match[1]
			var data map[string]interface{}

			err := json.Unmarshal(jsonContent, &data)
			if err != nil {
				continue
			}

			// Check if this is an array of JSON-LD objects
			if _, ok := data["@graph"]; ok {
				graph, ok := data["@graph"].([]interface{})
				if !ok {
					continue
				}

				for _, item := range graph {
					itemMap, ok := item.(map[string]interface{})
					if !ok {
						continue
					}

					if typeStr, ok := itemMap["@type"].(string); ok && typeStr == "Recipe" {
						jsonBytes, err := json.Marshal(itemMap)
						if err != nil {
							continue
						}

						var recipe JSONLDRecipe
						if err := json.Unmarshal(jsonBytes, &recipe); err == nil {
							return &recipe
						}
					}
				}
				continue
			}

			// Check if this is a recipe directly
			if typeStr, ok := data["@type"].(string); ok {
				if typeStr == "Recipe" {
					var recipe JSONLDRecipe
					if err := json.Unmarshal(jsonContent, &recipe); err == nil {
						return &recipe
					}
				}
			}
		}
	}
	return nil
}

// mapJSONLDToRecipe maps JSON-LD data to our Recipe struct
func mapJSONLDToRecipe(jsonLD *JSONLDRecipe, recipe *ExtractedRecipe) {
	recipe.Title = jsonLD.Name
	recipe.Description = jsonLD.Description
	recipe.PrepTime = jsonLD.PrepTime
	recipe.CookTime = jsonLD.CookTime
	recipe.TotalTime = jsonLD.TotalTime

	// Handle image URL (can be string or array)
	switch img := jsonLD.Image.(type) {
	case string:
		recipe.ImageURL = img
	case []interface{}:
		if len(img) > 0 {
			if imgStr, ok := img[0].(string); ok {
				recipe.ImageURL = imgStr
			}
		}
	case map[string]interface{}:
		if url, ok := img["url"].(string); ok {
			recipe.ImageURL = url
		}
	}

	// Handle author (can be string or object)
	switch author := jsonLD.Author.(type) {
	case string:
		recipe.Author = author
	case map[string]interface{}:
		if name, ok := author["name"].(string); ok {
			recipe.Author = name
		}
	}

	// Handle yield (can be string or array)
	switch yield := jsonLD.RecipeYield.(type) {
	case string:
		recipe.Yield = yield
	case []interface{}:
		if len(yield) > 0 {
			if yieldStr, ok := yield[0].(string); ok {
				recipe.Yield = yieldStr
			}
		}
	}

	// Handle ingredients (array of strings)
	switch ingredients := jsonLD.Ingredients.(type) {
	case []interface{}:
		for _, ing := range ingredients {
			if ingStr, ok := ing.(string); ok {
				recipe.Ingredients = append(recipe.Ingredients, strings.TrimSpace(ingStr))
			}
		}
	}

	// Handle instructions (array of strings or objects)
	switch instructions := jsonLD.Instructions.(type) {
	case []interface{}:
		for _, inst := range instructions {
			if instStr, ok := inst.(string); ok {
				recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(instStr))
			} else if instMap, ok := inst.(map[string]interface{}); ok {
				if text, ok := instMap["text"].(string); ok {
					recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(text))
				}
			}
		}
	case string:
		// Some sites just provide a single string with line breaks
		lines := strings.Split(instructions, "\n")
		for _, line := range lines {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				recipe.Instructions = append(recipe.Instructions, trimmed)
			}
		}
	}
}

// extractMicrodata extracts schema.org Recipe microdata
func extractMicrodata(doc *goquery.Document, recipe *ExtractedRecipe) {
	// Look for elements with itemtype="http://schema.org/Recipe"
	doc.Find("[itemtype='http://schema.org/Recipe'], [itemtype='https://schema.org/Recipe']").Each(func(_ int, s *goquery.Selection) {
		// Extract ingredients
		s.Find("[itemprop='recipeIngredient'], [itemprop='ingredients']").Each(func(_ int, ing *goquery.Selection) {
			text := strings.TrimSpace(ing.Text())
			if text != "" {
				recipe.Ingredients = append(recipe.Ingredients, text)
			}
		})

		// Extract instructions
		s.Find("[itemprop='recipeInstructions']").Each(func(_ int, inst *goquery.Selection) {
			// Check if instructions are in a list
			inst.Find("li").Each(func(_ int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if text != "" {
					recipe.Instructions = append(recipe.Instructions, text)
				}
			})

			// If no list items were found, use the text directly
			if len(recipe.Instructions) == 0 {
				text := strings.TrimSpace(inst.Text())
				if text != "" {
					// Split by lines or sentences if it's a block of text
					lines := strings.Split(text, "\n")
					for _, line := range lines {
						if trimmed := strings.TrimSpace(line); trimmed != "" {
							recipe.Instructions = append(recipe.Instructions, trimmed)
						}
					}
				}
			}
		})
	})
}

// extractIngredientsByPatterns looks for common patterns for ingredients
func extractIngredientsByPatterns(doc *goquery.Document, recipe *ExtractedRecipe) {
	// Common class names and patterns for ingredient lists
	selectors := []string{
		".ingredients", "#ingredients",
		"[class*='ingredient']", "ul[class*='ingredient']",
		".ing-list", ".recipe-ingredients",
		"[id*='ingredient']",
	}

	for _, selector := range selectors {
		// Try to find container
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			var items []string

			// Try to find list items within the container
			s.Find("li").Each(func(_ int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if text != "" {
					items = append(items, text)
				}
			})

			// If list items found, use them
			if len(items) > 0 {
				recipe.Ingredients = items
				return
			}

			// Otherwise, look for separate elements like divs or spans
			s.Find("div, span, p").Each(func(_ int, el *goquery.Selection) {
				text := strings.TrimSpace(el.Text())
				if text != "" && len(text) < 200 && !strings.Contains(strings.ToLower(text), "instruction") {
					items = append(items, text)
				}
			})

			// If items found in divs or spans, use them
			if len(items) > 0 {
				recipe.Ingredients = items
				return
			}

			// Last resort: use the text content of the container itself
			text := strings.TrimSpace(s.Text())
			if text != "" {
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					if trimmed := strings.TrimSpace(line); trimmed != "" && trimmed != "Ingredients" {
						recipe.Ingredients = append(recipe.Ingredients, trimmed)
					}
				}
			}
		})

		// If ingredients were found, break the loop
		if len(recipe.Ingredients) > 0 {
			break
		}
	}
}

// extractInstructionsByPatterns looks for common patterns for instructions
func extractInstructionsByPatterns(doc *goquery.Document, recipe *ExtractedRecipe) {
	// Common class names and patterns for instruction lists
	selectors := []string{
		".instructions", "#instructions",
		"[class*='instruction']", "ol[class*='instruction']",
		".steps", ".recipe-directions", ".recipe-steps",
		"[id*='direction']", "[class*='direction']",
		".method", "#method",
	}

	for _, selector := range selectors {
		// Try to find container
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			var items []string

			// Try to find list items within the container
			s.Find("li").Each(func(_ int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if text != "" {
					items = append(items, text)
				}
			})

			// If list items found, use them
			if len(items) > 0 {
				recipe.Instructions = items
				return
			}

			// Look for step headings or numbered paragraphs
			s.Find("h3, h4, p").Each(func(_ int, el *goquery.Selection) {
				text := strings.TrimSpace(el.Text())
				if text != "" && (strings.HasPrefix(text, "Step") ||
					(len(text) > 1 && text[0] >= '1' && text[0] <= '9' && text[1] == '.')) {
					items = append(items, text)
				}
			})

			// If step headings found, use them
			if len(items) > 0 {
				recipe.Instructions = items
				return
			}

			// Look for div or p elements that might contain steps
			s.Find("div, p").Each(func(_ int, el *goquery.Selection) {
				text := strings.TrimSpace(el.Text())
				if text != "" && len(text) > 20 && !strings.Contains(strings.ToLower(text), "ingredient") {
					items = append(items, text)
				}
			})

			// If items found in divs or paragraphs, use them
			if len(items) > 0 {
				recipe.Instructions = items
				return
			}

			// Last resort: use the text content of the container itself
			text := strings.TrimSpace(s.Text())
			if text != "" {
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					if trimmed := strings.TrimSpace(line); trimmed != "" &&
						trimmed != "Instructions" && trimmed != "Directions" && len(trimmed) > 10 {
						recipe.Instructions = append(recipe.Instructions, trimmed)
					}
				}
			}
		})

		// If instructions were found, break the loop
		if len(recipe.Instructions) > 0 {
			break
		}
	}
}

// extractUnstructuredIngredients tries to find ingredients in unstructured content
func extractUnstructuredIngredients(doc *goquery.Document, recipe *ExtractedRecipe) {
	// Look for paragraphs with ingredient-like content
	doc.Find("p").Each(func(_ int, p *goquery.Selection) {
		text := strings.TrimSpace(p.Text())

		// Check if the paragraph likely contains ingredients
		if strings.Contains(strings.ToLower(text), "cup") ||
			strings.Contains(strings.ToLower(text), "tablespoon") ||
			strings.Contains(strings.ToLower(text), "teaspoon") ||
			strings.Contains(strings.ToLower(text), "ounce") ||
			strings.Contains(strings.ToLower(text), "gram") {

			// Split by common separators
			items := splitTextByPattern(text, []string{",", "•", "\\*", ";", "\\n"})
			if len(items) > 1 {
				for _, item := range items {
					if trimmed := strings.TrimSpace(item); trimmed != "" {
						recipe.Ingredients = append(recipe.Ingredients, trimmed)
					}
				}
			}
		}
	})
}

// extractUnstructuredInstructions tries to find instructions in unstructured content
func extractUnstructuredInstructions(doc *goquery.Document, recipe *ExtractedRecipe) {
	// Look for paragraphs with instruction-like content
	doc.Find("p, div").Each(func(_ int, el *goquery.Selection) {
		text := strings.TrimSpace(el.Text())

		// Check if the element likely contains instructions
		lowerText := strings.ToLower(text)
		if (strings.Contains(lowerText, "mix") ||
			strings.Contains(lowerText, "stir") ||
			strings.Contains(lowerText, "heat") ||
			strings.Contains(lowerText, "cook") ||
			strings.Contains(lowerText, "bake")) &&
			len(text) > 50 && !strings.Contains(lowerText, "ingredient") {

			// Split by sentences or numbered patterns
			if strings.Contains(text, ".") {
				sentences := splitIntoSentences(text)
				for _, sentence := range sentences {
					if trimmed := strings.TrimSpace(sentence); trimmed != "" && len(trimmed) > 15 {
						recipe.Instructions = append(recipe.Instructions, trimmed)
					}
				}
			} else {
				// Try to split by numbered steps
				re := regexp.MustCompile(`(\d+\.\s+)`)
				steps := re.Split(text, -1)
				if len(steps) > 1 {
					for _, step := range steps {
						if trimmed := strings.TrimSpace(step); trimmed != "" && len(trimmed) > 15 {
							recipe.Instructions = append(recipe.Instructions, trimmed)
						}
					}
				} else {
					// Just use the whole paragraph
					recipe.Instructions = append(recipe.Instructions, text)
				}
			}
		}
	})
}

// splitTextByPattern splits text using a list of patterns
func splitTextByPattern(text string, patterns []string) []string {
	pattern := strings.Join(patterns, "|")
	re := regexp.MustCompile(pattern)
	return re.Split(text, -1)
}

// splitIntoSentences splits text into sentences
func splitIntoSentences(text string) []string {
	// Simple sentence splitter (not perfect but works for most cases)
	re := regexp.MustCompile(`[.!?]+\s+`)
	sentences := re.Split(text, -1)

	// Handle the last sentence which might end with period without space
	if !strings.HasSuffix(text, ". ") && strings.HasSuffix(text, ".") {
		lastPart := sentences[len(sentences)-1]
		if strings.HasSuffix(lastPart, ".") {
			sentences[len(sentences)-1] = lastPart[:len(lastPart)-1]
		}
	}

	return sentences
}
