/*
	Inputs finds takes a user's profile and calculates their daily macro-nutrient requirements. This can be fed into
	the Outputs package which generates the user's personalised recipe based on these requirements.

	The user's profile is stored in a struct containing the following variables:
		Weight - measured in kg
		Height - measured in cm
		Body Fat - dimensionless ratio (bewteen 0 and 1
		Activity Level - a modifier based on the user's daily exercise amount (between 1 and 2)
		Age
		Gender
		Regime - the user's desired regime: bulk, cut or maintain

	The user's macro-nutrient requirements are stored in a struct containing the following variables
		Carbs - the user's recommended daily carbohydrate intake, measured in g
		Protein - the user's recommended daily protein intake, measured in g
		Fat - the user's recommended daily fat intake, measured in g
		Calories - the user's recommended daily calorie intake, measured in kcal
*/
package inputs

import (
	"fmt"
	"math"
	"os/user"
)

const kilo2Pound = 2.20462262 	// kilogram to pound conversion factor
const pMassRatio = 1.5			// number of grams of protein required per lb of lean body mass
const fMassRatio = 0.45			// number of grams of fat required per lb of body mass
const pCalRatio = 4				// number of calories in each gram of protein
const fCalRatio = 9				// number of calories in each gram of fat
const cCalRatio = 4				// number of calories in each gram of carbohydrate
const bulkModifier = 1.1		// calorie modifier for bulk regime
const cutModifier = 0.9			// calorie modifier for cut regime

type Profile struct {
	Weight, Height, BodyFat, ActivityLevel, Age  float64
	Gender, Regime string
}

type Macros struct {
	Carbs, Protein, Fat, Calories float64
}

//leanBodyMass finds the user's lean body mass from their weight and body fat
func leanBodyMass(weight, bodyFat float64) float64 {
	return weight * (1-bodyFat)
}

//findCalories finds the user's recommended daily calorie intake based on the average of three well-known methodologies.
//the result is multiplied by a modifier to suit the user's desired regime.
func (u *Profile) findCalories() float64 {

	var harrisBenedict, mifflinStJeor float64

	switch u.Gender {
	case "male":
		harrisBenedict = 66 + (13.7 * u.Weight) + (5 * u.Height) - (6.76 * u.Age)
		mifflinStJeor = (9.99 * u.Weight) + (6.25 * u.Height) - (4.92 * u.Age) + 5
	case "female":
		harrisBenedict = 655 + (9.6 * u.Weight) + (1.8 * u.Height) - (4.7 * u.Age)
		mifflinStJeor = (9.99 * u.Weight) + (6.25 * u.Height) - (4.92 * u.Age) - 161
	}
	
	katchMcArdle := 370 + (leanBodyMass(u.Weight, u.BodyFat) * 21.6)

	calories := u.ActivityLevel * (harrisBenedict + mifflinStJeor + katchMcArdle) / 3

	switch u.Regime {
	case "bulk":
		return calories * bulkModifier
	case "maintain":
		return calories
	case "cut":
		return calories * cutModifier
	}
}

/*
	FindMacros outputs a Macros struct containing the user's recommended daily macro-nutrient intake.

	The protein requirement is calculated as the product of the user's lean body mass and a pre-defined ratio.
	The fat requirement is calculated as the product of the user's weight and a pre-defined ratio.
	The calorie requirement is calculated using the findCalories function.
	The carbohydrate requirement is the amount of carbohydrates required to complete the user's calorie intake after
	taking into account the calories gained from their fat and protein requirements.
*/
func (u *Profile) FindMacros() Macros {

	UserMacros := Macros{}
	
	UserMacros.Protein = pMassRatio * leanBodyMass(u.Weight, u.BodyFat) * kilo2Pound
	UserMacros.Fat = fMassRatio * u.Weight * kilo2Pound
	UserMacros.Calories = u.findCalories()
	UserMacros.Carbs = (UserMacros.Calories - pCalRatio*UserMacros.Protein - fCalRatio*UserMacros.Fat)/cCalRatio

	return UserMacros
}

