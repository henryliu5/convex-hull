import numpy as np
import matplotlib.pyplot as plt
import math
import sys

def plot_points(x,y):
    for i in range(0, len(x)):
        plt.plot(x[i], y[i], 'ro-')

def plot_points_pts(pts, color='ro-'):
    for i in range(0, len(pts)):
        plt.plot(pts[i][0], pts[i][1], color)

def plot_hull(x,y):
    x = x.copy()
    y = y.copy()
    
    x.append(x[0])
    y.append(y[0])
    for i in range(0, len(x)-1):
        plt.plot(x[i:i+2], y[i:i+2], 'bo-')
        
def plot_hull_pts(pts):
    x = list()
    y = list()
    for i in range(len(pts)):
        x.append(pts[i][0])
        y.append(pts[i][1])
    plot_hull(x,y)

def get_mean_pt(x,y):
    return [sum(x)/len(x), sum(y)/len(y)]

def verify_points(hull_x, hull_y, x, y, save_dir, plot_graph=True):
    
    hull_x = hull_x.copy()
    hull_y = hull_y.copy()
    
    hull_x.append(hull_x[0])
    hull_y.append(hull_y[0])
    
    total_failed = 0
    
    for i in range(len(x)):
        failed = 0
        for j in range(len(hull_x) - 1):
            base_hull = (hull_x[j], hull_y[j])
            next_hull = (hull_x[j+1], hull_y[j+1])
            
            
            vec_to_pos = np.array([x[i] - base_hull[0], y[i] - base_hull[1]])
            vec_to_hull = np.array([next_hull[0] - base_hull[0], next_hull[1] - base_hull[1]])
            
            #Not in polygon
            if (np.cross(vec_to_hull, vec_to_pos) <= 0):
                failed = 1
                break
        
        if (plot_graph):
            if (failed == 1):
                plt.plot(x[i], y[i], 'ro-')
            else:
                plt.plot(x[i], y[i], 'go-')
        total_failed += failed
        
    if (plot_graph):
        plot_hull(hull_x, hull_y)
        plt.savefig(save_dir)
    return failed == 0


def polar_coord_sort(origin, x_ary, y_ary):
    pts = list()
    for i in range(len(x_ary)):
        x,y = x_ary[i] - origin[0], y_ary[i] - origin[1]
        theta, r = math.atan(y/x), math.sqrt(x**2 + y**2) 
        if (y < 0 and x > 0):
            theta = theta + 2 * math.pi
        elif (x < 0):
            theta = theta + math.pi
        pts.append([theta, r, x_ary[i], y_ary[i]])
    
    pts = sorted(pts)
    x_pts = list()
    y_pts = list()
    for a in pts:
        x_pts.append(a[2])
        y_pts.append(a[3])
    return x_pts, y_pts

points_file_dir = sys.argv[1]
hull_file_dir = sys.argv[2]
save_dir = sys.argv[3]

points_x, points_y = list(), list()
points_file = open(points_file_dir)
for line in points_file:
    x,y = float(line.split(",")[0]), float(line.split(",")[1])
    points_x.append(x)
    points_y.append(y)


hull_x, hull_y = list(), list()
hull_file = open(hull_file_dir)
for line in hull_file:
    x,y = float(line.split(",")[0]), float(line.split(",")[1])
    hull_x.append(x)
    hull_y.append(y)

hull_x, hull_y = polar_coord_sort(get_mean_pt(hull_x,hull_y), hull_x, hull_y)
verify_points(hull_x, hull_y, points_x, points_y, save_dir, plot_graph=True)

