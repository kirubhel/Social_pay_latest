/**
 * Color extraction utility for analyzing images and extracting dominant colors
 */

export interface ColorInfo {
  hex: string;
  rgb: [number, number, number];
  hsl: [number, number, number];
  frequency: number;
}

export interface ColorPalette {
  primary: ColorInfo;
  secondary: ColorInfo;
  accent: ColorInfo;
  dominantColors: ColorInfo[];
}

/**
 * Convert RGB to HSL
 */
function rgbToHsl(r: number, g: number, b: number): [number, number, number] {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0;
  let s = 0;
  const l = (max + min) / 2;

  if (max === min) {
    h = s = 0; // achromatic
  } else {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return [Math.round(h * 360), Math.round(s * 100), Math.round(l * 100)];
}

/**
 * Convert RGB to HEX
 */
function rgbToHex(r: number, g: number, b: number): string {
  return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

/**
 * Get color distance between two RGB colors
 */
function colorDistance(rgb1: [number, number, number], rgb2: [number, number, number]): number {
  const [r1, g1, b1] = rgb1;
  const [r2, g2, b2] = rgb2;
  return Math.sqrt(Math.pow(r2 - r1, 2) + Math.pow(g2 - g1, 2) + Math.pow(b2 - b1, 2));
}

/**
 * Check if color is too dark or too light (grayscale)
 */
function isGrayscale(rgb: [number, number, number], threshold = 15): boolean {
  const [r, g, b] = rgb;
  return Math.abs(r - g) < threshold && Math.abs(g - b) < threshold && Math.abs(r - b) < threshold;
}

/**
 * Extract colors from an image
 */
export async function extractColorsFromImage(imageUrl: string): Promise<ColorPalette> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous';
    
    img.onload = () => {
      try {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        
        if (!ctx) {
          reject(new Error('Could not get canvas context'));
          return;
        }

        // Resize for performance while maintaining aspect ratio
        const maxSize = 200;
        const scale = Math.min(maxSize / img.width, maxSize / img.height);
        canvas.width = img.width * scale;
        canvas.height = img.height * scale;

        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
        
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        const pixels = imageData.data;
        
        // Count color frequencies
        const colorMap = new Map<string, { rgb: [number, number, number], count: number }>();
        
        for (let i = 0; i < pixels.length; i += 4) {
          const r = pixels[i];
          const g = pixels[i + 1];
          const b = pixels[i + 2];
          const a = pixels[i + 3];
          
          // Skip transparent pixels
          if (a < 128) continue;
          
          // Skip very light or very dark colors
          const brightness = (r + g + b) / 3;
          if (brightness < 20 || brightness > 235) continue;
          
          // Group similar colors together (reduce color precision)
          const groupedR = Math.round(r / 10) * 10;
          const groupedG = Math.round(g / 10) * 10;
          const groupedB = Math.round(b / 10) * 10;
          
          const rgb: [number, number, number] = [groupedR, groupedG, groupedB];
          const key = `${groupedR},${groupedG},${groupedB}`;
          
          if (colorMap.has(key)) {
            colorMap.get(key)!.count++;
          } else {
            colorMap.set(key, { rgb, count: 1 });
          }
        }
        
        // Convert to array and sort by frequency
        const colorArray = Array.from(colorMap.entries())
          .map(([key, data]) => ({
            hex: rgbToHex(...data.rgb),
            rgb: data.rgb,
            hsl: rgbToHsl(...data.rgb),
            frequency: data.count
          }))
          .filter(color => !isGrayscale(color.rgb))
          .sort((a, b) => b.frequency - a.frequency);
        
        if (colorArray.length === 0) {
          reject(new Error('No significant colors found in image'));
          return;
        }
        
        // Find primary and secondary colors
        const primary = colorArray[0];
        let secondary = colorArray[1];
        
        // Ensure secondary is sufficiently different from primary
        for (let i = 1; i < colorArray.length; i++) {
          if (colorDistance(primary.rgb, colorArray[i].rgb) > 50) {
            secondary = colorArray[i];
            break;
          }
        }
        
        // Find accent color (something complementary or vibrant)
        let accent = colorArray[2] || secondary;
        for (let i = 0; i < colorArray.length; i++) {
          const color = colorArray[i];
          const [h, s, l] = color.hsl;
          // Look for vibrant colors (high saturation)
          if (s > 60 && l > 20 && l < 80 && colorDistance(primary.rgb, color.rgb) > 30) {
            accent = color;
            break;
          }
        }
        
        const palette: ColorPalette = {
          primary,
          secondary: secondary || primary,
          accent,
          dominantColors: colorArray.slice(0, 8)
        };
        
        resolve(palette);
      } catch (error) {
        reject(error);
      }
    };
    
    img.onerror = () => {
      reject(new Error('Failed to load image'));
    };
    
    img.src = imageUrl;
  });
}

/**
 * Generate a complete color scheme from extracted colors
 */
export function generateColorScheme(palette: ColorPalette) {
  const { primary, secondary, accent } = palette;
  
  // Generate color variations (lighter and darker shades)
  const generateShades = (baseColor: ColorInfo) => {
    const [r, g, b] = baseColor.rgb;
    
    return {
      50: rgbToHex(
        Math.min(255, Math.round(r + (255 - r) * 0.9)),
        Math.min(255, Math.round(g + (255 - g) * 0.9)),
        Math.min(255, Math.round(b + (255 - b) * 0.9))
      ),
      100: rgbToHex(
        Math.min(255, Math.round(r + (255 - r) * 0.8)),
        Math.min(255, Math.round(g + (255 - g) * 0.8)),
        Math.min(255, Math.round(b + (255 - b) * 0.8))
      ),
      200: rgbToHex(
        Math.min(255, Math.round(r + (255 - r) * 0.6)),
        Math.min(255, Math.round(g + (255 - g) * 0.6)),
        Math.min(255, Math.round(b + (255 - b) * 0.6))
      ),
      300: rgbToHex(
        Math.min(255, Math.round(r + (255 - r) * 0.4)),
        Math.min(255, Math.round(g + (255 - g) * 0.4)),
        Math.min(255, Math.round(b + (255 - b) * 0.4))
      ),
      400: rgbToHex(
        Math.min(255, Math.round(r + (255 - r) * 0.2)),
        Math.min(255, Math.round(g + (255 - g) * 0.2)),
        Math.min(255, Math.round(b + (255 - b) * 0.2))
      ),
      500: baseColor.hex,
      600: rgbToHex(
        Math.round(r * 0.8),
        Math.round(g * 0.8),
        Math.round(b * 0.8)
      ),
      700: rgbToHex(
        Math.round(r * 0.6),
        Math.round(g * 0.6),
        Math.round(b * 0.6)
      ),
      800: rgbToHex(
        Math.round(r * 0.4),
        Math.round(g * 0.4),
        Math.round(b * 0.4)
      ),
      900: rgbToHex(
        Math.round(r * 0.2),
        Math.round(g * 0.2),
        Math.round(b * 0.2)
      )
    };
  };
  
  return {
    primary: generateShades(primary),
    secondary: generateShades(secondary),
    accent: generateShades(accent),
    gradients: {
      primary: `linear-gradient(135deg, ${primary.hex} 0%, ${secondary.hex} 100%)`,
      secondary: `linear-gradient(135deg, ${secondary.hex} 0%, ${accent.hex} 100%)`,
      accent: `linear-gradient(135deg, ${accent.hex} 0%, ${primary.hex} 100%)`
    }
  };
}

/**
 * Social Pay brand colors extracted from logo analysis
 * These can be used as fallback colors or default theme
 */
export const SOCIAL_PAY_BRAND_COLORS = {
  primary: {
    hex: '#22c55e',
    rgb: [34, 197, 94] as [number, number, number],
    hsl: [142, 70, 45] as [number, number, number],
    frequency: 0
  },
  secondary: {
    hex: '#fbbf24',
    rgb: [251, 191, 36] as [number, number, number],
    hsl: [44, 96, 56] as [number, number, number],
    frequency: 0
  },
  accent: {
    hex: '#06b6d4',
    rgb: [6, 182, 212] as [number, number, number],
    hsl: [189, 94, 43] as [number, number, number],
    frequency: 0
  }
};

/**
 * Apply extracted colors to CSS custom properties
 */
export function applyColorScheme(scheme: ReturnType<typeof generateColorScheme>) {
  const root = document.documentElement;
  
  // Apply primary colors
  Object.entries(scheme.primary).forEach(([shade, color]) => {
    root.style.setProperty(`--color-primary-${shade}`, color);
  });
  
  // Apply secondary colors
  Object.entries(scheme.secondary).forEach(([shade, color]) => {
    root.style.setProperty(`--color-secondary-${shade}`, color);
  });
  
  // Apply accent colors
  Object.entries(scheme.accent).forEach(([shade, color]) => {
    root.style.setProperty(`--color-accent-${shade}`, color);
  });
  
  // Apply gradients
  Object.entries(scheme.gradients).forEach(([name, gradient]) => {
    root.style.setProperty(`--gradient-${name}`, gradient);
  });
} 