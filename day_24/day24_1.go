// Álvaro Castellano Vela 2019/01/07

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Group struct {
	Units          int
	HitPoints      int
	Weaknesses     map[string]bool
	Immunities     map[string]bool
	Damage         int
	DamageType     string
	Initiative     int
	SystemType     string
	EffectivePower int
}

func (attacker Group) calculateDamageover(target Group) int {
	//fmt.Println("calculateDamageover")
	//fmt.Println(attacker)
	//fmt.Println(target)
	//fmt.Println("_________________")
	var multiplier int = 1
	if _, ok := target.Weaknesses[attacker.DamageType]; ok {
		multiplier = 2
	} else {
		if _, ok := target.Immunities[attacker.DamageType]; ok {
			multiplier = 0
		}
	}

	return attacker.EffectivePower * multiplier
}

type TargetGroup struct {
	SystemType     string
	EffectivePower int
	GroupID        int
	Initiative     int
}

type InitiativeGroup struct {
	SystemType string
	GroupID    int
	Initiative int
	OriginalID int
}

type TargetGroups []TargetGroup

func (x TargetGroups) Len() int      { return len(x) }
func (x TargetGroups) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x TargetGroups) Less(i, j int) bool {
	var result bool

	if x[i].EffectivePower < x[j].EffectivePower {
		result = true
	} else {
		if x[i].EffectivePower > x[j].EffectivePower {
			result = false
		} else {
			if x[i].Initiative < x[j].Initiative {
				result = true
			} else {
				result = false
			}
		}
	}
	return result
}

type InitiativeGroups []InitiativeGroup

func (x InitiativeGroups) Len() int      { return len(x) }
func (x InitiativeGroups) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x InitiativeGroups) Less(i, j int) bool {
	var result bool = false

	if x[i].Initiative < x[j].Initiative {
		result = true
	}
	return result
}

type ImmuneSystem []Group
type InfectionSystem []Group
type Systems map[string][]Group

func processFile(filename string) Systems {

	var immuneGroups ImmuneSystem
	var infectionGroups InfectionSystem

	systems := make(Systems)

	var groups [][]Group = make([][]Group, 2)
	groups[0] = immuneGroups
	groups[1] = infectionGroups

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	//Immune system
	scanner.Scan()
	for i := 0; i < 2; i++ {
		for scanner.Scan() {
			var group Group
			immuneGroupString := scanner.Text()
			if immuneGroupString == "" {
				break
			}
			re := regexp.MustCompile("^([[:digit:]]+) units each with ([[:digit:]]+) hit points \\(([^\\)]+)\\) with an attack that does ([[:digit:]]+) ([[:alpha:]]+) damage at initiative ([[:digit:]]+)$")
			match := re.FindAllStringSubmatch(immuneGroupString, -1)

			matchLen := len(match)
			var offset int = 0
			if matchLen == 0 {
				offset = 0
				re = regexp.MustCompile("^([[:digit:]]+) units each with ([[:digit:]]+) hit points with an attack that does ([[:digit:]]+) ([[:alpha:]]+) damage at initiative ([[:digit:]]+)$")
				match = re.FindAllStringSubmatch(immuneGroupString, -1)
			} else {
				offset = 1
			}

			group.Units, _ = strconv.Atoi(match[0][1])
			group.HitPoints, _ = strconv.Atoi(match[0][2])
			group.Damage, _ = strconv.Atoi(match[0][3+offset])
			group.DamageType = match[0][4+offset]
			group.Initiative, _ = strconv.Atoi(match[0][5+offset])

			if matchLen == 1 {
				weaknessesImmunities := match[0][3]
				weaknessesImmunities = strings.Replace(weaknessesImmunities, ",", "", -1)
				weaknessesImmunities = strings.Replace(weaknessesImmunities, ";", "", -1)
				weaknessesImmunities = strings.Replace(weaknessesImmunities, "to", "", -1)
				weaknessesImmunities = strings.Replace(weaknessesImmunities, "  ", " ", -1)
				weaknessesImmunitiesSlice := strings.Split(weaknessesImmunities, " ")

				group.Weaknesses = make(map[string]bool)
				group.Immunities = make(map[string]bool)

				var currenType string
				for _, typeString := range weaknessesImmunitiesSlice {
					if typeString == "weak" || typeString == "immune" {
						currenType = typeString
					} else {
						if currenType == "weak" {
							group.Weaknesses[typeString] = true
						} else {
							group.Immunities[typeString] = true
						}
					}
				}
			}
			if i == 0 {
				group.SystemType = "Immune"
			} else {
				group.SystemType = "Infection"
			}
			group.EffectivePower = group.Units * group.Damage
			groups[i] = append(groups[i], group)
		}
		scanner.Scan()
	}
	systems["Immune"] = groups[0]
	systems["Infection"] = groups[1]
	return systems
}

func getAttackGroups(systems Systems) TargetGroups {
	var targetGroups TargetGroups

	for _, system := range systems {
		for id, group := range system {
			if group.Units > 0 {
				targetgroup := TargetGroup{SystemType: group.SystemType, EffectivePower: group.EffectivePower, GroupID: id, Initiative: group.Initiative}
				targetGroups = append(targetGroups, targetgroup)
			}
		}
	}
	sort.Sort(sort.Reverse(targetGroups))
	fmt.Println(targetGroups)

	return targetGroups
}

func fight(systems Systems) int {

	// 1- Target selection
	targetGroups := getAttackGroups(systems)
	InfectionAttackChoosal := make(map[int]int)
	ImmuneAtackChoosal := make(map[int]int)

	InfectionAttackChoosen := make(map[int]int)
	ImmuneAtackChoosen := make(map[int]int)

	var attackChoosal *map[int]int
	var attackChoosen *map[int]int
	var initiativeGroups InitiativeGroups

	for attackerID, attacker := range targetGroups {
		var systemToAttack string
		fmt.Println("attacker.SystemType: ", attacker.SystemType)
		if attacker.SystemType == "Immune" {
			systemToAttack = "Infection"
			attackChoosal = &InfectionAttackChoosal
			attackChoosen = &ImmuneAtackChoosen
		} else {
			systemToAttack = "Immune"
			attackChoosal = &ImmuneAtackChoosal
			attackChoosen = &InfectionAttackChoosen
		}
		attakerGroup := systems[attacker.SystemType][attacker.GroupID]
		var maxDamage int = 0
		var maxEffectivePower int = 0
		var maxInitiative int = 0
		var idToAttack int = -1
		for targetID, target := range systems[systemToAttack] {
			if _, ok := (*attackChoosen)[targetID]; !ok {
				damage := attakerGroup.calculateDamageover(target)
				if damage > maxDamage {
					maxDamage = damage
					maxEffectivePower = target.EffectivePower
					maxInitiative = target.Initiative
					idToAttack = targetID
				} else {
					if damage == maxDamage {
						if target.EffectivePower > maxEffectivePower {
							maxDamage = damage
							maxEffectivePower = target.EffectivePower
							maxInitiative = target.Initiative
							idToAttack = targetID
						} else {
							if target.Initiative > maxInitiative {
								maxDamage = damage
								maxEffectivePower = target.EffectivePower
								maxInitiative = target.Initiative
								idToAttack = targetID
							}
						}
					}
				}
			}
		}
		if idToAttack > -1 {
			(*attackChoosal)[attackerID] = idToAttack
			(*attackChoosen)[idToAttack] = attackerID
			fmt.Println()
			fmt.Println("attacker")
			fmt.Println(systems[attacker.SystemType][attacker.GroupID])
			fmt.Println("choses")
			fmt.Println(systems[systemToAttack][idToAttack])
			fmt.Println()
			initiativeGroups = append(initiativeGroups, InitiativeGroup{SystemType: attacker.SystemType, GroupID: attacker.GroupID, Initiative: attacker.Initiative, OriginalID: attackerID})
		}
	}

	sort.Sort(sort.Reverse(initiativeGroups))
	fmt.Println("Attack")
	fmt.Println(initiativeGroups)
	fmt.Println()
	fmt.Println()
	//Time to attack

	for _, initiativeGroup := range initiativeGroups {
		var attackerPtr, targetPtr *Group
		var attackerOrdered TargetGroup = targetGroups[initiativeGroup.OriginalID]
		var systemToAttack string

		if attackerOrdered.SystemType == "Immune" {
			systemToAttack = "Infection"
			attackChoosal = &InfectionAttackChoosal
		} else {
			systemToAttack = "Immune"
			attackChoosal = &ImmuneAtackChoosal
		}

		attackerPtr = &systems[initiativeGroup.SystemType][initiativeGroup.GroupID]
		targetPtr = &systems[systemToAttack][(*attackChoosal)[initiativeGroup.OriginalID]]

		fmt.Println("Attacker")
		fmt.Println(attackerPtr)
		fmt.Println("EffectivePower:", (*attackerPtr).EffectivePower)
		fmt.Println("Attacks")
		fmt.Println(targetPtr)
		damage := (*attackerPtr).calculateDamageover(*targetPtr)
		fmt.Printf("Damage will be %d\n", damage)
		fmt.Printf("HitPoints %d\n", (*targetPtr).HitPoints)
		killedUnits := damage / (*targetPtr).HitPoints
		fmt.Printf("Killed units %d\n", killedUnits)
		(*targetPtr).Units -= killedUnits
		(*targetPtr).EffectivePower = (*targetPtr).Units * (*targetPtr).Damage
	}

	return 0
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	systems := processFile(filename)

	result := fight(systems)

	fmt.Printf("Final Resuls %d\n", result)
}
