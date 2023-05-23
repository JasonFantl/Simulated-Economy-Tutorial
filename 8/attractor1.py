import math
import numpy as np

import matplotlib
import matplotlib.pyplot as plt

def aizawa(x, y, z, a=0.95, b=0.7, c=0.65, d=3.5, e=0.25, f=0.1):
    x_prime = (z - b) * x - d * y
    y_prime = d * x + (z - b) * y
    z_prime = c + a * z - (z ** 3) / 3. - (x ** 2 + y ** 2) * (1 + e * z) + f * z * x ** 3
    return x_prime, y_prime, z_prime

# how many steps we propagate
nb_iters = int(1e7) #1e7
dt = 0.01

# rather than draw the attractor directly, let's
# instead use a histogram to determine where we
# spend most of our time
h, w = 1600, 1600
hist = np.zeros((h, w), dtype=int)

x_min = -2
x_max = 2

y_min = x_min * (h / w) + 0.6
y_max = x_max * (h / w) + 0.4

print('x_min: ', x_min)
print('x_max: ', x_max)
print('y_min: ', y_min)
print('y_max: ', y_max)

x, y, z = 0.1, 0, 0

for i in range(nb_iters):
    x_prime, y_prime, z_prime = aizawa(x, y, z)

    x += x_prime * dt
    y += y_prime * dt
    z += z_prime * dt

    # allow particles to settle
    if nb_iters < int(5e1):
        continue

    x_i = int( (x - x_min) * w / (x_max - x_min) )
    # we are only plotting the x and z axes here
    y_i = int( (z - y_min) * h / (y_max - y_min) )
    if (x_i >= 0 and x_i < w \
        and y_i >= 0 and y_i < h):
        # remember that matrix is row by columns
        hist[y_i, x_i] += 1

im = np.ones((h, w, 3), dtype=float)
sens = 3e-4
color = (200, 200, 200)
for i in range(h):
    for j in range(w):
        val = hist[i, j]
        r = math.exp(-sens * val * color[0])
        g = math.exp(-sens * val * color[1])
        b = math.exp(-sens * val * color[2])
        im[i, j, :] = r, g, b

plt.imsave('aizawa.png', im, dpi=600, origin='lower')

plt.axis('off')
plt.imshow(im)
plt.show()