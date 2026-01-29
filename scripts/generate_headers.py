
import os

HEADERS_DIR = "web/public/headers"

# Common Definitions
DEFS = """
  <defs>
     <linearGradient id="divine_gradient" x1="0" y1="0" x2="100%" y2="100%" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#06b6d4"/> <!-- Cyan-500 -->
      <stop offset="50%" stop-color="#8b5cf6"/> <!-- Violet-500 -->
      <stop offset="100%" stop-color="#eab308"/> <!-- Yellow-500 (Gold) -->
    </linearGradient>
    <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="4" result="blur"/>
      <feComposite in="SourceGraphic" in2="blur" operator="over"/>
    </filter>
  </defs>
"""

STYLE = """
  <style>
    .title { font-family: system-ui, -apple-system, sans-serif; font-weight: 800; font-size: 80px; fill: url(#divine_gradient); letter-spacing: -1px; }
    .subtitle { font-family: system-ui, -apple-system, sans-serif; font-weight: 500; font-size: 24px; fill: #94a3b8; letter-spacing: 6px; text-transform: uppercase; }
    .icon-stroke { stroke: url(#divine_gradient); stroke-width: 6; fill: none; stroke-linecap: round; stroke-linejoin: round; filter: url(#glow); }
    .icon-fill { fill: url(#divine_gradient); filter: url(#glow); }
  </style>
"""

MODULES = [
    {
        "name": "memos",
        "parts": ["ME", "MOS"],
        "parts_zh": ["闪念", "笔记"],
        "icon": """
            <!-- Divine Spark Card Icon -->
            <!-- Card Outline -->
            <rect class="icon-stroke" x="40" y="30" width="120" height="140" rx="20" />
            
            <!-- Inner Divine Spark (4-pointed star) -->
            <path class="icon-fill" d="M100 55 L115 90 L150 100 L115 110 L100 145 L85 110 L50 100 L85 90 Z" />
            
            <!-- Orbit / Flow lines -->
            <path class="icon-stroke" d="M30 100 Q 100 180 170 100" stroke-width="3" opacity="0.5" stroke-dasharray="4 4"/>
            <circle class="icon-fill" cx="100" cy="100" r="8" fill="#fff"/>
        """
    },
    {
        "name": "ai",
        "parts": ["AI", "CHAT"], # Assistant is too long/hard to split nicely. "AI CHAT" is better? Or "ASST" "TANT"? Let's use "AI" "ASST" or just "ASSIS" "TANT". User said "2 chars" so short is better. Let's try proper split.
        # Actually user said "2 zi" (2 characters). Chinese is 4 chars total usually.
        # Let's stick to the meaning. "Assistant" -> "ASST" "ANT"?
        # Let's use "ASSIS" "TANT"
        "parts": ["ASSIS", "TANT"],
        "parts_zh": ["智能", "助理"],
        "icon": """
            <path class="icon-stroke" d="M100 30 C 60 30 30 60 30 100 C 30 140 60 170 100 170 C 140 170 170 140 170 100 C 170 60 140 30 100 30 Z" opacity="0.5"/>
            <path class="icon-stroke" d="M100 100 L 70 70 M100 100 L 130 70 M100 100 L 70 130 M100 100 L 130 130" />
            <circle class="icon-fill" cx="100" cy="100" r="15" />
            <circle class="icon-fill" cx="70" cy="70" r="6" />
            <circle class="icon-fill" cx="130" cy="70" r="6" />
            <circle class="icon-fill" cx="70" cy="130" r="6" />
            <circle class="icon-fill" cx="130" cy="130" r="6" />
            <circle class="icon-fill" cx="100" cy="30" r="4" />
            <circle class="icon-fill" cx="30" cy="100" r="4" />
            <circle class="icon-fill" cx="170" cy="100" r="4" />
        """
    },
    {
        "name": "schedule",
        "parts": ["SCHE", "DULE"],
        "parts_zh": ["日程", "管理"],
        "icon": """
            <rect class="icon-stroke" x="30" y="40" width="140" height="130" rx="16" />
            <path class="icon-stroke" d="M130 20 L130 60 M70 20 L70 60" />
            <path class="icon-stroke" d="M30 80 L170 80" />
            <circle class="icon-fill" cx="100" cy="125" r="25" opacity="0.8" />
        """
    },
    {
        "name": "review",
        "parts": ["RE", "VIEW"],
        "parts_zh": ["每日", "回顾"],
        "icon": """
             <path class="icon-stroke" d="M100 170 A 70 70 0 1 1 170 100" />    
             <path class="icon-stroke" d="M170 100 L170 70 L195 95" stroke-linecap="square" />
             <circle class="icon-fill" cx="100" cy="100" r="15" />
             <path class="icon-stroke" d="M100 60 L100 75 M100 125 L100 140 M60 100 L75 100 M125 100 L140 100" stroke-width="4"/>
        """
    },
    {
        "name": "knowledge",
        "parts": ["KNOW", "LEDGE"],
        "parts_zh": ["知识", "图谱"],
        "icon": """
            <circle class="icon-stroke" cx="50" cy="150" r="20" />
            <circle class="icon-stroke" cx="150" cy="150" r="20" />
            <circle class="icon-stroke" cx="100" cy="50" r="20" />
            <line class="icon-stroke" x1="50" y1="130" x2="100" y2="70" />
            <line class="icon-stroke" x1="150" y1="130" x2="100" y2="70" />
            <line class="icon-stroke" x1="70" y1="150" x2="130" y2="150" />
            <circle class="icon-fill" cx="100" cy="100" r="10" />
            <line class="icon-stroke" x1="50" y1="130" x2="100" y2="100" opacity="0.5"/>
            <line class="icon-stroke" x1="150" y1="130" x2="100" y2="100" opacity="0.5"/>
            <line class="icon-stroke" x1="100" y1="70" x2="100" y2="100" opacity="0.5"/>
        """
    },
    {
        "name": "explore",
        "parts": ["EX", "PLORE"],
        "parts_zh": ["探索", "发现"],
        "icon": """
            <circle class="icon-stroke" cx="100" cy="100" r="70" />
            <path class="icon-stroke" d="M100 30 L100 170 M30 100 L170 100" opacity="0.5" />
            <path class="icon-fill" d="M100 60 L115 100 L100 140 L85 100 Z" />
        """
    },
    {
        "name": "files",
        "parts": ["RES", "OURCES"],
        "parts_zh": ["资源", "附件"],
        "icon": """
             <path class="icon-stroke" d="M30 160 L50 60 C52 50 60 40 70 40 L170 40 L190 160 Z" /> 
             <path class="icon-stroke" d="M30 160 L190 160" />
             <path class="icon-stroke" d="M70 40 L90 20 L150 20 L170 40"/>
             <rect class="icon-fill" x="90" y="80" width="40" height="50" rx="4" />
        """
    },
    {
        "name": "inbox",
        "parts": ["IN", "BOX"],
        "parts_zh": ["消息", "通知"],
        "icon": """
            <path class="icon-stroke" d="M100 30 C 70 30 50 50 50 90 C 50 140 30 150 30 150 L 170 150 C 170 150 150 140 150 90 C 150 50 130 30 100 30" />
            <path class="icon-stroke" d="M85 150 C 85 160 92 170 100 170 C 108 170 115 160 115 150" />
            <circle class="icon-fill" cx="150" cy="50" r="10" />
        """
    }
]


def generate_svg(module, is_zh=False):
    parts = module.get('parts_zh') if is_zh else module.get('parts')
    part1, part2 = parts[0], parts[1]
    
    # Layout: [Part1]  [Icon]  [Part2]
    # Centered in VIEW_W
    
    VIEW_W = 1200
    VIEW_H = 140
    
    # Text metrics
    title_size = 72 
    
    # Icon settings
    icon_scale = 0.7 
    icon_w_unscaled = 200 
    icon_w = icon_w_unscaled * icon_scale # ~140
    
    gap = 25 
    
    # Y Alignment
    text_y = "50%"
    icon_y = 0 
    
    center_x = VIEW_W / 2

    if is_zh:
        # --- Chinese Layout: Split [Part1] [Icon] [Part2] ---
        
        return f"""<svg width="{VIEW_W}" height="{VIEW_H}" viewBox="0 0 {VIEW_W} {VIEW_H}" fill="none" xmlns="http://www.w3.org/2000/svg">
{DEFS}
{STYLE}
  <!-- Content Group -->
  
  <!-- Icon centered at center_x -->
  <g transform="translate({center_x - (icon_w/2)}, {icon_y}) scale({icon_scale})">
      {module['icon']}
  </g>
  
  <!-- Part 1: Aligned END (Right) to the Left of Icon -->
  <g transform="translate({center_x - (icon_w/2) - gap}, 0)">
       <text class="title" x="0" y="{text_y}" style="font-size: {title_size}px; text-anchor: end; dominant-baseline: central; alignment-baseline: central;">{part1}</text>
  </g>

  <!-- Part 2: Aligned START (Left) to the Right of Icon -->
  <g transform="translate({center_x + (icon_w/2) + gap}, 0)">
       <text class="title" x="0" y="{text_y}" style="font-size: {title_size}px; text-anchor: start; dominant-baseline: central; alignment-baseline: central;">{part2}</text>
  </g>
</svg>"""

    else:
        # --- English Layout: [Title] [Icon] (Icon on the Right) ---
        
        full_title = "".join(parts)
        
        # Estimate Widths for Centering
        # English: Increased estimate to 60px/char for better centering accuracy with current font size (72px)
        # Font aspect ratio ~0.8 * 72 = 57.6, so 60 is safe/good.
        char_w = 60 
        title_w = len(full_title) * char_w
        
        # Total Group Width
        total_w = title_w + gap + icon_w
        
        # Start X for the group
        start_x = (VIEW_W - total_w) / 2
        
        return f"""<svg width="{VIEW_W}" height="{VIEW_H}" viewBox="0 0 {VIEW_W} {VIEW_H}" fill="none" xmlns="http://www.w3.org/2000/svg">
{DEFS}
{STYLE}
  <!-- Content Group starting at start_x -->
  
  <!-- Title: Aligned START (Left) -->
  <g transform="translate({start_x}, 0)">
       <text class="title" x="0" y="{text_y}" style="font-size: {title_size}px; text-anchor: start; dominant-baseline: central; alignment-baseline: central;">{full_title}</text>
  </g>
  
  <!-- Icon: To the Right of Title -->
  <!-- Icon X = start_x + title_w + gap -->
  <!-- Note: We use estimated title_w for positioning icon? 
       Using text-anchor start for title is good, but we need precise end of text to place icon.
       SVG doesn't auto-flow. 
       Better approach: Center the text and icon relative to center_x, but that requires knowing exact text width.
       Alternative: Use text-anchor="end" for title at center? No, user wants [Title] [Icon].
       
       Let's stick to estimated centering but align Icon based on estimate.
       OR, we can align [Title] [Icon] group to center.
       
       Let's anchor Title at {start_x}. 
       Icon at {start_x + title_w + gap}.
  -->
  <g transform="translate({start_x + title_w + gap}, {icon_y}) scale({icon_scale})">
      {module['icon']}
  </g>

</svg>"""

def main():
    if not os.path.exists(HEADERS_DIR):
        os.makedirs(HEADERS_DIR)
        
    for module in MODULES:
        # Generate English version
        filename_en = f"header-{module['name']}.svg"
        filepath_en = os.path.join(HEADERS_DIR, filename_en)
        content_en = generate_svg(module, is_zh=False)
        with open(filepath_en, 'w') as f:
            f.write(content_en)
        print(f"Generated {filepath_en}")
        
        # Generate Chinese version
        filename_zh = f"header-{module['name']}-zh.svg"
        filepath_zh = os.path.join(HEADERS_DIR, filename_zh)
        content_zh = generate_svg(module, is_zh=True)
        with open(filepath_zh, 'w') as f:
            f.write(content_zh)
        print(f"Generated {filepath_zh}")

if __name__ == "__main__":
    main()
