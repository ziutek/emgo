package vec3d

type Matrix struct {
	v1, v2, v3 Vec
}

type Vec struct {
	x, y, z float32
}

func (v Vec) Dot(a Vec) float32 {
	return v.x*a.x + v.y*a.y + v.z*a.z
}

func (v Vec) XYZ() (x, y, z float32) {
	return v.x, v.y, v.z
}

func (v Vec) Add(a Vec) (r Vec) {
	r.x = v.x + a.x
	r.y = v.y + a.y
	r.z = v.z + a.z
	return
}

func (v Vec) V3() (vx, vy, vz Vec) {
	vx.x = v.x
	vy.y = v.y
	vz.z = v.z
	return
}
