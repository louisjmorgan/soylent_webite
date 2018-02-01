/*
	Outputs generates a soylent recipe based on the user's macros.
*/
package outputs

import "github.com/louisjmorgan/soylent_website/lib/inputs"

const oatBase = 200
const oatMod = 150
const calUpper = 3400
const calLower = 1600
const fibreRatio = 14
const salt = 4.0
const multivitamin = 1.8
const choline = 1.0
const potassiumReq = 2500

type Macros inputs.Macros

type Ingredient struct {
	Protein, Carbs, Fat, Fibre, Potassium float64
}

type Recipe struct {
	Oats, Whey, Maltodextrin, Oil, Psyllium, Salt, Multivitamin, Choline, Potassium float64
}

/*
	GenerateRecipe determines the quantities of each ingredient required to fulfill the user's macronutrient
	requirements.

	The oat flour amount is chosen first as it contains the widest range of nutrients and therefore serves as the base
	of the recipe. The amount is scaled to the user's daily calorie intake from a base amount. The base amount,
	modification amount and upper and lower calorie bounds are chosen semi-arbitrarily to provide a good balance between
	oats and maltodextrin for the sources of carbohydrate; limiting the amount of maltodextrin without prescribing so
	much oats as to cause overdosing of certain micronutrients (especially iron) in the cases of users with
	exceptionally high calorie requirements.

	Once the amount of oat flour has been chosen, the whey protein can then be set to fulfill the remaining protein
	requirement of the user. This works because the remaining ingredients have insignifcant protein contents.

	The amount of oil is set in the same way, as only the oat flour and whey protein contain significant amounts of fat.

	The psyllium husk amount is then calculated based on a ratio between required calorie intake and require dietary
	fibre, taking into account the fibre contributed by the oat flour and the whey protein.

	The maltodextrin amount is then chosen to fulfil the user's remaining carbohydrate requirement, taking into account
	the carbohydrates contributed by the oat flour, whey protein and psyllium husks.

	Finally, the potassium gluconate amount is set to fulfill the user's potassium requirements, taking into account
	the potassium contribution of the oat flour and the multivitamin powder.

	The remaining ingredients have fixed values as they are contributing trace micronutrients whose exact value is not
	important provided it's within certain bounds.
*/
func (m *Macros) GenerateRecipe() Recipe {

	Recipe := Recipe{}

	var (
		Oats, Whey, Maltodextrin, Oil, Psyllium, Salt, Multivitamin, Choline, Potassium Ingredient
	)

	Recipe.Salt = salt
	Recipe.Multivitamin = multivitamin
	Recipe.Choline = choline
	Recipe.Oats = oatBase + oatMod*(m.Calories-calLower)/(calUpper-calLower)
	Recipe.Whey = (m.Protein - Recipe.Oats*Oats.Protein) / Whey.Protein
	Recipe.Psyllium = ((m.Calories/1000)*fibreRatio - Recipe.Whey*Whey.Fibre - Recipe.Oats*Oats.Fibre) / Psyllium.Fibre
	Recipe.Oil = (m.Fat - Recipe.Oats*Oats.Fat - Recipe.Whey*Whey.Fat)
	Recipe.Maltodextrin = (m.Carbs - Recipe.Oats*Oats.Carbs - Recipe.Whey*Whey.Carbs - Recipe.Psyllium - Psyllium.Carbs)
	Recipe.Potassium = (potassiumReq - Recipe.Oats*Oats.Potassium - Recipe.Multivitamin*Multivitamin.Potassium)

	return Recipe
}
