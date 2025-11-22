# RendIm
A path tracer based on Peter Shirley's [Ray](http://www.amazon.com/gp/product/B01B5AODD8/ref=as_li_tl) [Tracing](https://www.amazon.com/gp/product/B01CO7PQ8C/ref=series_rw_dp_sw) [Minibooks](https://www.amazon.com/gp/product/B01DN58P8C/ref=series_rw_dp_sw)

## Features

### Rendering Engine
- **Path tracing** with configurable samples per pixel
- **Multi-threaded rendering** with bucket-based parallel processing
- **Real-time preview** via WebSocket streaming
- **BVH acceleration** (Bounding Volume Hierarchy) for faster ray-object intersection

### Materials
- **Lambertian** (diffuse) surfaces
- **Metal** with adjustable fuzziness
- **Dielectric** (glass) with refraction and Schlick approximation
- **Diffuse light** emitters
- **Isotropic** scattering for volumes

### Textures
- **Constant** color textures
- **Checker** patterns
- **Perlin noise** with turbulence
- **Image mapping** (JPEG support)

### Geometry
- **Spheres** (static and moving)
- **Axis-aligned rectangles** (XY, XZ, YZ planes)
- **Boxes**
- **Volumes** (constant density medium for fog/smoke)
- **Transformations**: translation, Y-axis rotation
- **Normal flipping** for inside-out surfaces

### Camera
- **Positionable** camera with look-from/look-at
- **Adjustable field of view**
- **Depth of field** (defocus blur) with configurable aperture
- **Motion blur** with shutter time interval

Image from the cover of the first book:

![alt text](https://github.com/MiroslavGatsanoga/RendIm/blob/master/out.png)

Rendering the image with 100 samples per pixel took **49m 4s** on machine with Intel Core i5-4210H CPU 2.90 GHz.

Image from the cover of the second book:

![alt text](https://github.com/MiroslavGatsanoga/RendIm/blob/master/out2.png)

Rendering the image with 10000 samples per pixel took **36h 43m 16s** on machine with Intel Core i5-4210H CPU 2.90 GHz.
