# go-graphics
Basic raytracing and rasterization in go + SDL (https://github.com/veandco/go-sdl2).

`cd ray && SDL_RENDER_DRIVER=software go run .`

![first-ray-trace](ray/images/BasicRaytracing.png)


### Basic Concepts
* Canvas - a rectangular array of colored pixels.
* Coordinates 
  * Screen coordinates - origin is the top left
  * Canvas coordinates - origin is center
* Scene - the set of objects to render.
* Scene Coordinates
  * Y is up, X and Z are horizontal


## Ray Tracing

### Coordinate Conversion
* Camera Position = "O" = (0,0,0)
* Camera Orientation = Looks towards the positive Z-axis (with Y+ up and X+ to the right)
* Viewport = (Vw x Vh), Frontal to the camera (perpendicular to Z+ axis) and at distance d.

Since the viewport center matches the canvas center and the orientation is the same, in order to convert from canvas coords to viewport coords, we just need to scale:

Vx = Cx * (Vw/Cw)
Vy = Cy * (Vy/Cy)

Additionally, ever point on the viewport has distance "d", so z = "d" by defintion:

Vz = d

### Tracing Rays
In reality, light comes from a light source, but it's more efficient for us to do the tracing in reverse.
We'll consider a ray that originates at the camera, goes through a point in the viewport, and then hits an object in the scene.

### Equations
* Ray Equation -> P = O + t(V-O) (where t is any real number representing the distance along the ray)
  * Another way to express this is P = O +t*D (where D is the direction of the ray)

### Sphere Equation
* Based on "math", we can transform the equation for points on a sphere into a more useful form:
  * <P-C,P-C> = r^2 (Where the first part is the dot product of the vector from point to center)

### Combination
* We want to find P values that satisfy both equations (points that lie on the ray and on the sphere)
* We can combine both equations and simplify using trig to end up with:
  * t^2<D,D> + t(2<CO,D>) + <CO,CO> - r^2 = 0
  * This has the form of a quadratic equation, which has either 0, 1, or 2 solutions.
    * 0 solutions -> ray never hits the sphere.
    * 1 solution -> ray hits tangentially.
    * 2 solutions -> ray hits the sphere once on each side.
* Considering t
  * t < 0 represents points "behind the eye"
  * t = 1 represents points on the viewport
  * t > 1 represents points behind the viewpoint (these are the ones we care about)

### Lessons Learned
* The canvas and computer screen coordinate system conversion can be tricky. I initially used the following simple formulas:
    * Sx = Cw/2 + Cx
    * Sy = Cy/2 - Cy
* Mistake 1: I missed the subtraction sign in the second one. My scene rendered, but it was upside down.
* Mistake 2: I fixed the minus sign, then I hit an array out of bounds. You actually have to subtract 1.
* It's hard to visualize, but it actually makes sense. The loop I was using went from -Y to Y-1 (due to less than comparison).
* Once I flipped to subtraction, then I had a mismatch on my bounds and thus the overflow.
