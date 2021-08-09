package main

import "github.com/gobwas/glob"
import "fmt"

func main() {
	var g glob.Glob

	// Examples

	//1.Case : Match group level
	g = glob.MustCompile("*grp_123*")

	fmt.Println("Matching group level")
	fmt.Println("*grp_123*", "/RightsSamples/grp_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123**/.SinglePrimitive"))
	fmt.Println("*grp_123*", "/RightsSamples/grp_123/comp_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123/comp_123**/.SinglePrimitive"))
	fmt.Println("*grp_123*", "/RightsSamples/grp_123/comp_123/loc_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123/comp_123/loc_123**/.SinglePrimitive"))
	fmt.Println("*grp_123*", "/RightsSamples/grp_123/comp_123/loc_123/cust_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123/comp_123/loc_123/cust_123**/.SinglePrimitive"))

	//1.Case : Match location level
	g = glob.MustCompile("*grp_123/comp_123/loc_123*")

	fmt.Println("Matching location level")

	fmt.Println("*grp_123/comp_123/loc_123*", "/RightsSamples/grp_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123**/.SinglePrimitive"))
	fmt.Println("*grp_123/comp_123/loc_123*", "/RightsSamples/grp_123/comp_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123/comp_123**/.SinglePrimitive"))
	fmt.Println("*grp_123/comp_123/loc_123*", "/RightsSamples/grp_123/comp_123/loc_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123/comp_123/loc_123**/.SinglePrimitive"))
	fmt.Println("*grp_123/comp_123/loc_123*", "/RightsSamples/grp_123/comp_123/loc_123/cust_123**/.SinglePrimitive", g.Match("/RightsSamples/grp_123/comp_123/loc_123/cust_123**/.SinglePrimitive"))
}
