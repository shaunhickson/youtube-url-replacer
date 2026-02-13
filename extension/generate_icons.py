from PIL import Image, ImageDraw, ImageFont
import os

def create_icon(size):
    # Create image with transparent background
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # Background: Deep Purple Rounded Square
    bg_color = (103, 58, 183, 255) # Deep Purple
    rect_coords = [0, 0, size, size]
    draw.rounded_rectangle(rect_coords, radius=size//5, fill=bg_color)

    # Speech Bubble (White)
    bubble_w = size * 0.7
    bubble_h = size * 0.5
    bubble_x = (size - bubble_w) / 2
    bubble_y = (size - bubble_h) / 2 - (size * 0.05)
    
    draw.rounded_rectangle(
        [bubble_x, bubble_y, bubble_x + bubble_w, bubble_y + bubble_h], 
        radius=size//8, 
        fill=(255, 255, 255, 255)
    )
    
    # Bubble Tail
    tail_coords = [
        (bubble_x + bubble_w * 0.2, bubble_y + bubble_h - 1),
        (bubble_x + bubble_w * 0.2, bubble_y + bubble_h + size * 0.15),
        (bubble_x + bubble_w * 0.5, bubble_y + bubble_h - 1)
    ]
    draw.polygon(tail_coords, fill=(255, 255, 255, 255))

    # Play Button (Purple) inside the bubble
    play_size = bubble_h * 0.5
    play_x = bubble_x + (bubble_w - play_size) / 2 + (play_size * 0.1) # nudge right for visual center
    play_y = bubble_y + (bubble_h - play_size) / 2
    
    p1 = (play_x, play_y)
    p2 = (play_x, play_y + play_size)
    p3 = (play_x + play_size * 0.866, play_y + play_size / 2)
    
    draw.polygon([p1, p2, p3], fill=bg_color)

    return img

sizes = [16, 48, 128]
base_path = '/Users/sph/src/youtube-url-replacer/extension/public/icons'

if not os.path.exists(base_path):
    os.makedirs(base_path)

for size in sizes:
    img = create_icon(size)
    img.save(f'{base_path}/icon{size}.png', 'PNG')
    print(f"Generated icon{size}.png")
