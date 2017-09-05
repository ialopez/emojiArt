package emojiart

import (
	"math"
	"sort"
)

type Tree struct {
	R, G, B  float64 //average rgb
	URLIndex int     //references index in emojiURLPath
	Left     *Tree
	Right    *Tree
}

/*interfaces and methods needed by sort
 */

type ByRed []Emoji
type ByGreen []Emoji
type ByBlue []Emoji

func (a ByRed) Len() int   { return len(a) }
func (a ByGreen) Len() int { return len(a) }
func (a ByBlue) Len() int  { return len(a) }
func (a ByRed) Swap(i, j int) {
	a[i].R, a[j].R = a[j].R, a[i].R
	a[i].G, a[j].G = a[j].G, a[i].G
	a[i].B, a[j].B = a[j].B, a[i].B
	a[i].URLIndex, a[j].URLIndex = a[j].URLIndex, a[i].URLIndex
}
func (a ByGreen) Swap(i, j int) {
	a[i].R, a[j].R = a[j].R, a[i].R
	a[i].G, a[j].G = a[j].G, a[i].G
	a[i].B, a[j].B = a[j].B, a[i].B
	a[i].URLIndex, a[j].URLIndex = a[j].URLIndex, a[i].URLIndex
}

func (a ByBlue) Swap(i, j int) {
	a[i].R, a[j].R = a[j].R, a[i].R
	a[i].G, a[j].G = a[j].G, a[i].G
	a[i].B, a[j].B = a[j].B, a[i].B
	a[i].URLIndex, a[j].URLIndex = a[j].URLIndex, a[i].URLIndex
}

func (a ByRed) Less(i, j int) bool   { return a[i].R < a[j].R }
func (a ByGreen) Less(i, j int) bool { return a[i].G < a[j].G }
func (a ByBlue) Less(i, j int) bool  { return a[i].B < a[j].B }

var kdTree map[string]*Tree

/*builds a balanced k-d tree
 */
func BuildTree(platform string, low, high, key int) *Tree {
	if (high < low) || (low > high) {
		//base case
		return nil
	} else if low == high {
		//base case
		node1 := Tree{
			R:        emojiDictAvg[platform][low].R,
			G:        emojiDictAvg[platform][low].G,
			B:        emojiDictAvg[platform][low].B,
			URLIndex: emojiDictAvg[platform][low].URLIndex,
		}
		return &node1
	} else {
		//recursive step
		//sort colors by current key
		switch key {
		case 0:
			sort.Sort(ByRed(emojiDictAvg[platform][low:high]))
		case 1:
			sort.Sort(ByGreen(emojiDictAvg[platform][low:high]))
		case 2:
			sort.Sort(ByBlue(emojiDictAvg[platform][low:high]))
		}

		//get median
		medianIndex := (low + high) / 2
		medianColor := emojiDictAvg[platform][medianIndex]
		parent := Tree{
			R:        medianColor.R,
			G:        medianColor.G,
			B:        medianColor.B,
			URLIndex: medianColor.URLIndex,
		}

		//create Tree struct for median and build its left and right children recursively
		key = (key + 1) % 3
		left := BuildTree(platform, low, medianIndex-1, key)
		right := BuildTree(platform, medianIndex+1, high, key)

		parent.Left = left
		parent.Right = right

		return &parent
	}
}

/*uses kdTree to find nearestNeighbor in log(n) time
 */
func nearestNeighbor(r, g, b float64, key int, node *Tree) (smallestDistance float64, URLIndex int) {
	//base case
	if node == nil {
		smallestDistance, URLIndex = -1, -1
		return //smallestDistance and URLIndex cannot be negative, check for -1 in recursive step
	} else {
		//smallest distance is set to value from recursive call, this value is the squared euclidean distance
		//take the square root if you want the actual distance

		var isLeftOfPlane bool //used to tell if point r, g ,b is left of the plane created by the current node
		nextKey := (key + 1) % 3
		switch key {
		case 0:
			if r < (*node).R {
				smallestDistance, URLIndex = nearestNeighbor(r, g, b, nextKey, (*node).Left)
				isLeftOfPlane = true
			} else {
				smallestDistance, URLIndex = nearestNeighbor(r, g, b, nextKey, (*node).Right)
				isLeftOfPlane = false
			}
		case 1:
			if g < (*node).G {
				smallestDistance, URLIndex = nearestNeighbor(r, g, b, nextKey, (*node).Left)
				isLeftOfPlane = true
			} else {
				smallestDistance, URLIndex = nearestNeighbor(r, g, b, nextKey, (*node).Right)
				isLeftOfPlane = false
			}
		case 2:
			if b < (*node).B {
				smallestDistance, URLIndex = nearestNeighbor(r, g, b, nextKey, (*node).Left)
				isLeftOfPlane = true
			} else {
				smallestDistance, URLIndex = nearestNeighbor(r, g, b, nextKey, (*node).Right)
				isLeftOfPlane = false
			}
		}

		distanceToThisNode := math.Pow(r-(*node).R, 2) + math.Pow(g-(*node).G, 2) + math.Pow(b-(*node).B, 2)

		//check if value in smallest distance is -1, if it is then we
		//set smallest distance to distance from the current node
		if smallestDistance == -1 {
			smallestDistance, URLIndex = distanceToThisNode, (*node).URLIndex
		} else if smallestDistance > distanceToThisNode {
			//else if smallest distance is bigger than distance from current node set to distance from this node
			smallestDistance, URLIndex = distanceToThisNode, (*node).URLIndex
		}

		//calculate if the sphere with center coordinates r, g, b and radius equal to smallestDistance intersects with
		//the xy-plane, xz-plane, or yz-plane(depending on the value of key) that contains the point R, G, B
		//from the current node
		//if it does calculate the nearest neighbor on the other side of the plane and save it to potential smallest
		var potentialSmallest float64
		var potentialURL int
		switch key {
		case 0:
			//yz plane
			if (isLeftOfPlane) && (r+math.Sqrt(smallestDistance) > node.R) {
				//sphere intersects
				potentialSmallest, potentialURL = nearestNeighbor(r, g, b, nextKey, (*node).Right)
			} else if (!isLeftOfPlane) && (r-math.Sqrt(smallestDistance) < node.R) {
				//sphere intersects
				potentialSmallest, potentialURL = nearestNeighbor(r, g, b, nextKey, (*node).Left)
			} else {
				potentialSmallest = -1
			}
		case 1:
			//xz plane
			if (isLeftOfPlane) && (g+math.Sqrt(smallestDistance) > node.G) {
				//sphere intersects
				potentialSmallest, potentialURL = nearestNeighbor(r, g, b, nextKey, (*node).Right)
			} else if (!isLeftOfPlane) && (g-math.Sqrt(smallestDistance) < node.G) {
				//sphere intersects
				potentialSmallest, potentialURL = nearestNeighbor(r, g, b, nextKey, (*node).Left)
			} else {
				//sphere doesnt intersect
				potentialSmallest = -1
			}
		case 2:
			//xy plane
			if (isLeftOfPlane) && (b+math.Sqrt(smallestDistance) > node.B) {
				//sphere intersects
				potentialSmallest, potentialURL = nearestNeighbor(r, g, b, nextKey, (*node).Right)
			} else if (!isLeftOfPlane) && (b-math.Sqrt(smallestDistance) < node.B) {
				//sphere intersects
				potentialSmallest, potentialURL = nearestNeighbor(r, g, b, nextKey, (*node).Left)
			} else {
				//sphere doesnt intersect
				potentialSmallest = -1
			}
		}

		//check if potential smallest is valid, if it is then check it is smaller than smallestDistance
		if potentialSmallest != -1 && potentialSmallest < smallestDistance {
			smallestDistance, URLIndex = potentialSmallest, potentialURL
		}
		return
	}
}
