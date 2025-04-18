package parsing

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/liushuangls/go-anthropic/v2"
)

type RecipeCollection struct {
	Recipes []Recipe `json:"recipes"`
}

// Recipe represents a single recipe with all its details
type Recipe struct {
	Name         string       `json:"name"`
	PrepTime     string       `json:"prep_time,omitempty"`
	CookTime     string       `json:"cook_time,omitempty"`
	TotalTime    string       `json:"total_time,omitempty"`
	Servings     int          `json:"servings,omitempty"`
	Cuisine      []string     `json:"cuisine,omitempty"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions,omitempty"`
	Notes        []string     `json:"notes,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
}

// Ingredient represents a single ingredient with its details
type Ingredient struct {
	Amount   *float64 `json:"amount"` // Pointer so it can be null
	Unit     string   `json:"unit,omitempty"`
	Name     string   `json:"name"`
	Notes    string   `json:"notes,omitempty"`
	Category string   `json:"category,omitempty"`
}

// RecipeTextToJSON converts recipe text to a JSON structure using Claude API
func RecipeTextToJsonString(text []byte) string {
	// Create a schema string based on our struct definitions
	schemaPrompt := `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Recipe Collection Schema",
  "description": "A schema for standardizing recipes for a shopping list planner",
  "type": "object",
  "required": ["recipes"],
  "properties": {
    "recipes": {
      "type": "array",
      "description": "An array of recipes",
      "items": {
        "type": "object",
        "required": ["name", "ingredients"],
        "properties": {
          "name": { "type": "string", "description": "The name of the recipe" },
          "prep_time": { "type": "string", "description": "Preparation time for the recipe" },
          "cook_time": { "type": "string", "description": "Cooking time for the recipe" },
          "total_time": { "type": "string", "description": "Total time to prepare the recipe" },
          "servings": { "type": "integer", "description": "Number of servings the recipe makes" },
          "cuisine": { "type": "array", "items": { "type": "string" }, "description": "Cuisines associated with the recipe" },
          "ingredients": {
            "type": "array",
            "description": "List of ingredients needed for the recipe",
            "items": {
              "type": "object",
              "required": ["name"],
              "properties": {
                "amount": { "type": ["number", "null"], "description": "Quantity of the ingredient" },
                "unit": { "type": "string", "description": "Unit of measurement for the ingredient" },
                "name": { "type": "string", "description": "Name of the ingredient" },
                "notes": { "type": "string", "description": "Additional notes about the ingredient" },
                "category": { "type": "string", "description": "Category of the ingredient" }
              }
            }
          },
          "instructions": { "type": "array", "description": "Step-by-step instructions for preparing the recipe", "items": { "type": "string" } },
          "notes": { "type": "array", "description": "Additional notes or tips for the recipe", "items": { "type": "string" } },
          "tags": { "type": "array", "description": "Tags for categorizing the recipe", "items": { "type": "string" } }
        }
      }
    }
  }
}`

	// Create the prompt with the schema and recipe text
	prompt := fmt.Sprintf(`
Extract the recipe information from the following text and format it according to this schema:
%s

Here's the recipe text to extract information from:
%s

Please return ONLY valid JSON that follows the schema with no additional text or explanation.

Note: Use tbsp or tsp rather than tablespoon or teaspoon. A correct item in the schema will look something like amount = 1, unit = cup, name = white whine, though this would appear as white wine 1 cup in the given text.
If an ingredient is listed as 'or', such as 4 lemons or 4 limes, select just the one which seem most appropriate for the meal.
Prefer to keep units simplified and standardized. Rather than 'large cloves', just clove. This is for a shopping list, so the aim is to make the data easily composable with other sources of data.
If a name is plural, give it as singular. DO NOT use plural forms.
In cases where there is no unit, such as tomatoes, we can ignore it. The unit does not denote size, that is apart of the name.
Given a range, choose a reasonable amount instead, avoid giving two options.
Maintain enough specificity so that someone can find the necessary items at the store. They want to know that it's greek yogurt rather than just yogurt, for instance.
Some recipes have multiple parts, like a dressing or sauce to go with it, that is why we specify a list of recipes in the json, they are all the subparts that make up the recipe as a whole.
Leave salt and pepper out, they are not needed to shop for most of the time.
Specify paprika as either smoked or sweet. Clove is not a unit (garlic)
`, schemaPrompt, string(text))

	// Initialize the Claude client
	client := anthropic.NewClient(os.Getenv("ANTHROPIC_KEY"))

	// Make the API call to Claude
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Dot5HaikuLatest,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(prompt),
		},
		MaxTokens: 2000,
	})

	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return ""
	}

	// Get the generated text
	generatedJSON := resp.Content[0].GetText()
	return generatedJSON
	// // Parse the generated JSON
	// var collection RecipeCollection
	// err = json.Unmarshal([]byte(generatedJSON), &collection)
	// if err != nil {
	// 	fmt.Printf("Error parsing generated JSON: %v\n", err)
	// 	return nil
	// }

	// return &collection
}

func RecipeIngredients(collection *RecipeCollection) string {
	if collection == nil {
		return ""
	}

	var list string
	for _, recipe := range collection.Recipes {
		for _, ingredient := range recipe.Ingredients {
			var amount string
			if ingredient.Amount != nil {
				amount = strconv.FormatFloat(*ingredient.Amount, 'f', -1, 64) + " "
			}
			list = fmt.Sprintf("%s%s %v%s\n", list, ingredient.Name, amount, ingredient.Unit) //FIXME use buffer
		}
	}
	return list
}
