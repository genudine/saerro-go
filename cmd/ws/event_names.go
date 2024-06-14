package main

import (
	"fmt"

	"github.com/genudine/saerro-go/util"
)

var experienceIDs = []int{
	2, 3, 4, 5, 6, 7, 34, 51, 53, 55, 57, 86, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
	100, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 201, 233, 293,
	294, 302, 303, 353, 354, 355, 438, 439, 503, 505, 579, 581, 584, 653, 656, 674, 675,
}

func getEventNames() []string {
	events := util.Map(experienceIDs, func(i int) string {
		return fmt.Sprintf("GainExperience_experience_id_%d", i)
	})
	events = append(events, "Death", "VehicleDestroy")

	return events
}
