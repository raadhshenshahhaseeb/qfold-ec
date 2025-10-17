## Definitions

### Affine Transformation
"Affine" describes a geometric transformation that is essentially a linear transformation (like rotation, scaling, or shearing) followed by a translation (a shift in position).

It doesn't have to keep the origin fixed. It performs a linear transformation and then shifts the entire space. This is why it's often described as a linear transformation plus a translation. The formula is: `new_vector = M * old_vector + b`, where `b` is the translation vector.

#### Affine Space
An affine space is a vector space (with points and vectors) but with no special origin. 

Think of a blank sheet of paper. It's an affine space. You can draw vectors (arrows) between any two points, and you can talk about directions and parallel lines. But until you pick a point and call it "`(0,0)`", there's no inherent "center" to the space.

Once you designate an origin, the affine space becomes a `vector space` because now every point can be described by a vector starting from that origin. In a way, an affine space is a vector space that has "forgotten" where its origin is.

##### Affine Coordinates `(x,y)` 
This is the standard, intuitive way to represent a point on a 2D plane. You have an x-value and a y-value. Addition and Doubling ops in ECC require modular inversion.

### Linear Transformation
A linear transformation always keeps the origin `(0,0)` fixed. Think of rotating an object around its center or scaling it from its center. The center point doesn't move. All linear transformations can be represented by a simple matrix multiplication: `new_vector = M * old_vector`.

### Analogy for Affine and Linear Transformation
Stretching or rotating the sheet while keeping the center point pinned down at the origin is a linear transformation.

Stretching, rotating and then sliding the sheet to a different spot is an affine transformation.

When we rotate the sheet about a point `p != (0,0)`, we are really doing 3 moves:
1) slide the sheet so `p goes to the origin`
2) do a plain rotation/scale (that part is linear)
3) slide back so `origin goes to p`

the first/last slide (translation) makes the whole thing affine, not linear. Linear maps can't slide.
