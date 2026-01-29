import sys
import os
from PIL import Image

def process_avatar(input_path, output_path):
    if not os.path.exists(input_path):
        print(f"Error: {input_path} not found")
        return

    img = Image.open(input_path).convert("RGBA")
    
    # Process using a simpler approach
    # Chroma Key: Remove pixels where Green is the dominant channel
    # This is more robust than fixed thresholds
    data = img.getdata()
    new_data = []
    
    for item in data:
        r, g, b, a = item
        # If green is very high and higher than others (Chroma Key)
        if g > 150 and g > r * 1.2 and g > b * 1.2:
            new_data.append((r, g, b, 0))
        # Also handle white background just in case
        elif r > 245 and g > 245 and b > 245:
            new_data.append((r, g, b, 0))
        else:
            new_data.append(item)
            
    img.putdata(new_data)
    img.save(output_path, "WEBP", quality=95)
    print(f"Success: Saved to {output_path}")

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python3 process_avatar.py <input> <output>")
    else:
        process_avatar(sys.argv[1], sys.argv[2])
