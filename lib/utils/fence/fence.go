// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package fence

import (
	"math"
	"reflect"
)

type DoubleFence struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//判断一个点是否在多边形内
func IsPtInPoly(point DoubleFence, pts []DoubleFence) (boundOrVertex bool) {

	var N = len(pts)
	boundOrVertex = true   //如果点位于多边形的顶点或边上，也算做点在多边形内，直接返回true
	var intersectCount = 0 //cross points count of x
	var precision = 2e-10  //浮点类型计算时候与0比较时候的容差
	var p1, p2 DoubleFence //neighbour bound vertices
	var p = point          //当前点

	p1 = pts[0] //left vertex
	for i := 0; i <= N; i++ { //check all rays
		if reflect.DeepEqual(p, p1) {
			return //p is an vertex
		}
		p2 = pts[i%N]

		//right vertex
		if p.X < math.Min(p1.X, p2.X) || p.X > math.Max(p1.X, p2.X) { //ray is outside of our interests
			p1 = p2
			continue //next ray left point
		}
		if p.X > math.Min(p1.X, p2.X) && p.X < math.Max(p1.X, p2.X) { //ray is crossing over by the algorithm (common part of)
			if p.Y <= math.Max(p1.Y, p2.Y) { //x is before of ray
				if p1.X == p2.X && p.Y >= math.Min(p1.Y, p2.Y) { //overlies on a horizontal ray
					return
				}
				if p1.Y == p2.Y { //ray is vertical
					if p1.Y == p.Y { //overlies on a vertical ray
						return
					} else { //before ray
						intersectCount++
					}
				} else { //cross point on the left side

					var xinters = (p.X-p1.X)*(p2.Y-p1.Y)/(p2.X-p1.X) + p1.Y //cross point of y
					if math.Abs(p.Y-xinters) < precision {                  //overlies on a ray
						return
					}

					if p.Y < xinters { //before ray
						intersectCount++
					}
				}
			}
		} else {                            //special case when ray is crossing through the vertex
			if p.X == p2.X && p.Y <= p2.Y { //p crossing over p2

				var p3 = pts[(i+1)%N]                                           //next vertex
				if p.X >= math.Min(p1.X, p3.X) && p.X <= math.Max(p1.X, p3.X) { //p.x lies between p1.x & p3.x
					intersectCount++
				} else {
					intersectCount += 2
				}
			}
		}
		p1 = p2 //next ray left point
	}

	if intersectCount%2 == 0 { //偶数在多边形外
		boundOrVertex = false
		return
	} //奇数在多边形内
	return

}
