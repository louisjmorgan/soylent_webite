//Finds a user's daily macro-nutrient requirements based on their profile
package find_macros

import (
	"fmt"
	"math"
	"os/user"
)

const kilo2pound = 2.20462262 // kilogram to pound conversion factor

//Data structure containing all the required user profile information
type Profile struct {
	weight float64			//in kg
	height float64			//in cm
	body_fat float64		//a dimensionless ratio between 0-1
	activity_level float64	//a number between 1-2
	age float64 			//a number
	gender string 			//"male" or "female"
	regime string			//"bulk", "cut" or "maintain"
}

//Data structure containing macro-nutrient information
type Macros struct {
	carbs, protein, fat, calories float64
}

//Initializing variables for the user's profile and macro-nutrients
var (
	user_profile = Profile{}
	user_macros = Macros{}
)

//Finds the user's daily calorie requirement based on an average of three different methodologies
func (u *Profile) avg_calories() float64 {

	var harris_benedict, mifflin_stjeor float64

	switch u.gender {
	case "male":
		harris_benedict = 66 + (13.7 * u.weight) + (5 * u.height) - (6.76 * u.age)
		mifflin_stjeor = (9.99 * u.weight) + (6.25 * u.height) - (4.92 * u.age) + 5
	case "female":
		harris_benedict = 655 + (9.6 * u.weight) + (1.8 * u.height) - (4.7 * u.age)
		mifflin_stjeor = (9.99 * u.weight) + (6.25 * u.height) - (4.92 * u.age) - 161
	}

	lean_body_mass := u.weight*(1-u.body_fat)
	katch_mcardle := 370 + (lean_body_mass*21.6)

	return (harris_benedict + mifflin_stjeor + katch_mcardle)/3
}

//
func (u *Profile) find_macros() Macros {

	weight_pounds := u.weight*kilo2pound
	lean_body_mass := weight_pounds*(1-u.body_fat)

	user_macros.protein = 1.5*lean_body_mass
	calories_protein := 4*user_macros.protein

	user_macros.fat = weight_pounds*0.45
	calories_fat := 9*user_macros.fat

	switch u.regime {
	case "bulk":
		user_macros.calories = u.avg_calories()*1.1
	case "maintain":
		user_macros.calories = u.avg_calories()
	case "cut":
		user_macros.calories = u.avg_calories()*0.9
	}

	user_macros.carbs = (user_macros.calories - calories_fat - calories_protein)/4

	return user_macros
}