'use client'

import { useState, useEffect } from 'react';
import Image from 'next/image';
import { 
  extractColorsFromImage, 
  generateColorScheme, 
  applyColorScheme, 
  SOCIAL_PAY_BRAND_COLORS,
  type ColorPalette 
} from '@/lib/color-extractor';

export default function TestCSSPage() {
  const [extractedColors, setExtractedColors] = useState<ColorPalette | null>(null);
  const [isExtracting, setIsExtracting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [colorScheme, setColorScheme] = useState<any>(null);

  const extractColors = async () => {
    setIsExtracting(true);
    setError(null);
    
    try {
      const palette = await extractColorsFromImage('/logo.png');
      setExtractedColors(palette);
      
      const scheme = generateColorScheme(palette);
      setColorScheme(scheme);
      
      // Apply the extracted colors to the page
      applyColorScheme(scheme);
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to extract colors');
      console.error('Color extraction error:', err);
    } finally {
      setIsExtracting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold mb-4 gradient-text-primary">
            Social Pay Color Extractor
          </h1>
          <p className="text-gray-600 text-lg">
            Extract and analyze colors from the Social Pay logo
          </p>
        </div>

        {/* Logo Display */}
        <div className="max-w-md mx-auto mb-8">
          <div className="bg-white rounded-2xl shadow-lg p-8 text-center">
            <h2 className="text-xl font-semibold mb-4 text-gray-800">Logo Analysis</h2>
            <div className="bg-gray-50 rounded-xl p-6 mb-6">
              <Image
                src="/logo.png"
                alt="Social Pay Logo"
                width={200}
                height={120}
                className="object-contain mx-auto"
                priority
              />
            </div>
            <button
              onClick={extractColors}
              disabled={isExtracting}
              className="w-full bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white font-semibold py-3 px-6 rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
            >
              {isExtracting ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-3"></div>
                  <span>Extracting Colors...</span>
                </div>
              ) : (
                'Extract Colors from Logo'
              )}
            </button>
          </div>
        </div>

        {/* Error Display */}
        {error && (
          <div className="max-w-2xl mx-auto mb-8">
            <div className="bg-red-50 border border-red-200 rounded-xl p-4">
              <p className="text-red-600 font-medium">Error: {error}</p>
            </div>
          </div>
        )}

        {/* Extracted Colors Display */}
        {extractedColors && (
          <div className="max-w-6xl mx-auto space-y-8">
            {/* Primary Colors */}
            <div className="bg-white rounded-2xl shadow-lg p-8">
              <h2 className="text-2xl font-bold mb-6 text-gray-800">Extracted Brand Colors</h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Primary Color */}
                <div className="text-center">
                  <div 
                    className="w-full h-24 rounded-xl shadow-lg mb-4 border-4 border-white"
                    style={{ backgroundColor: extractedColors.primary.hex }}
                  ></div>
                  <h3 className="font-semibold text-gray-800 mb-2">Primary Color</h3>
                  <div className="space-y-1 text-sm text-gray-600">
                    <p className="font-mono bg-gray-100 px-2 py-1 rounded">{extractedColors.primary.hex}</p>
                    <p>RGB({extractedColors.primary.rgb.join(', ')})</p>
                    <p>HSL({extractedColors.primary.hsl.join(', ')})</p>
                    <p>Frequency: {extractedColors.primary.frequency}</p>
                  </div>
                </div>

                {/* Secondary Color */}
                <div className="text-center">
                  <div 
                    className="w-full h-24 rounded-xl shadow-lg mb-4 border-4 border-white"
                    style={{ backgroundColor: extractedColors.secondary.hex }}
                  ></div>
                  <h3 className="font-semibold text-gray-800 mb-2">Secondary Color</h3>
                  <div className="space-y-1 text-sm text-gray-600">
                    <p className="font-mono bg-gray-100 px-2 py-1 rounded">{extractedColors.secondary.hex}</p>
                    <p>RGB({extractedColors.secondary.rgb.join(', ')})</p>
                    <p>HSL({extractedColors.secondary.hsl.join(', ')})</p>
                    <p>Frequency: {extractedColors.secondary.frequency}</p>
                  </div>
                </div>

                {/* Accent Color */}
                <div className="text-center">
                  <div 
                    className="w-full h-24 rounded-xl shadow-lg mb-4 border-4 border-white"
                    style={{ backgroundColor: extractedColors.accent.hex }}
                  ></div>
                  <h3 className="font-semibold text-gray-800 mb-2">Accent Color</h3>
                  <div className="space-y-1 text-sm text-gray-600">
                    <p className="font-mono bg-gray-100 px-2 py-1 rounded">{extractedColors.accent.hex}</p>
                    <p>RGB({extractedColors.accent.rgb.join(', ')})</p>
                    <p>HSL({extractedColors.accent.hsl.join(', ')})</p>
                    <p>Frequency: {extractedColors.accent.frequency}</p>
                  </div>
                </div>
              </div>
            </div>

            {/* Dominant Colors Palette */}
            <div className="bg-white rounded-2xl shadow-lg p-8">
              <h2 className="text-2xl font-bold mb-6 text-gray-800">Full Color Palette</h2>
              <div className="grid grid-cols-4 md:grid-cols-8 gap-4">
                {extractedColors.dominantColors.map((color, index) => (
                  <div key={index} className="text-center">
                    <div 
                      className="w-full h-16 rounded-lg shadow-md mb-2"
                      style={{ backgroundColor: color.hex }}
                    ></div>
                    <div className="text-xs space-y-1">
                      <p className="font-mono text-gray-700">{color.hex}</p>
                      <p className="text-gray-500">{color.frequency}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Color Scheme Generator */}
            {colorScheme && (
              <div className="bg-white rounded-2xl shadow-lg p-8">
                <h2 className="text-2xl font-bold mb-6 text-gray-800">Generated Color Scheme</h2>
                
                {/* Primary Shades */}
                <div className="mb-8">
                  <h3 className="text-lg font-semibold mb-4 text-gray-700">Primary Color Shades</h3>
                  <div className="grid grid-cols-9 gap-2">
                                         {Object.entries(colorScheme.primary).map(([shade, color]) => (
                       <div key={shade} className="text-center">
                         <div 
                           className="w-full h-12 rounded-lg shadow-sm"
                           style={{ backgroundColor: color as string }}
                         ></div>
                         <p className="text-xs mt-1 text-gray-600">{shade}</p>
                         <p className="text-xs font-mono text-gray-500">{color as string}</p>
                       </div>
                     ))}
                  </div>
                </div>

                {/* Secondary Shades */}
                <div className="mb-8">
                  <h3 className="text-lg font-semibold mb-4 text-gray-700">Secondary Color Shades</h3>
                  <div className="grid grid-cols-9 gap-2">
                                         {Object.entries(colorScheme.secondary).map(([shade, color]) => (
                       <div key={shade} className="text-center">
                         <div 
                           className="w-full h-12 rounded-lg shadow-sm"
                           style={{ backgroundColor: color as string }}
                         ></div>
                         <p className="text-xs mt-1 text-gray-600">{shade}</p>
                         <p className="text-xs font-mono text-gray-500">{color as string}</p>
                       </div>
                     ))}
                  </div>
                </div>

                {/* Gradients */}
                <div>
                  <h3 className="text-lg font-semibold mb-4 text-gray-700">Generated Gradients</h3>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                         {Object.entries(colorScheme.gradients).map(([name, gradient]) => (
                       <div key={name} className="text-center">
                         <div 
                           className="w-full h-20 rounded-xl shadow-lg mb-2"
                           style={{ background: gradient as string }}
                         ></div>
                         <h4 className="font-semibold text-gray-700 capitalize">{name} Gradient</h4>
                         <p className="text-xs font-mono text-gray-500 bg-gray-100 p-2 rounded mt-1">
                           {gradient as string}
                         </p>
                       </div>
                     ))}
                  </div>
                </div>
              </div>
            )}

            {/* Current vs New Design Comparison */}
            <div className="bg-white rounded-2xl shadow-lg p-8">
              <h2 className="text-2xl font-bold mb-6 text-gray-800">Design Examples</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                {/* Current Design */}
                <div>
                  <h3 className="text-lg font-semibold mb-4 text-gray-700">Current Brand Colors</h3>
                  <div className="space-y-4">
                    <button className="w-full bg-brand-green-500 hover:bg-brand-green-600 text-white font-semibold py-3 px-6 rounded-xl shadow-lg hover:shadow-xl transition-all duration-200">
                      Current Primary Button
                    </button>
                    <button className="w-full bg-brand-gold-400 hover:bg-brand-gold-500 text-white font-semibold py-3 px-6 rounded-xl shadow-lg hover:shadow-xl transition-all duration-200">
                      Current Secondary Button
                    </button>
                    <div className="bg-gradient-to-r from-brand-green-500 to-brand-gold-400 text-white p-4 rounded-xl">
                      <p className="font-semibold">Current Gradient Background</p>
                    </div>
                  </div>
                </div>

                {/* Extracted Design */}
                {extractedColors && (
                  <div>
                    <h3 className="text-lg font-semibold mb-4 text-gray-700">Extracted Colors</h3>
                    <div className="space-y-4">
                      <button 
                        className="w-full text-white font-semibold py-3 px-6 rounded-xl shadow-lg hover:shadow-xl transition-all duration-200"
                        style={{ backgroundColor: extractedColors.primary.hex }}
                      >
                        Extracted Primary Button
                      </button>
                      <button 
                        className="w-full text-white font-semibold py-3 px-6 rounded-xl shadow-lg hover:shadow-xl transition-all duration-200"
                        style={{ backgroundColor: extractedColors.secondary.hex }}
                      >
                        Extracted Secondary Button
                      </button>
                      <div 
                        className="text-white p-4 rounded-xl"
                        style={{ 
                          background: `linear-gradient(135deg, ${extractedColors.primary.hex} 0%, ${extractedColors.secondary.hex} 100%)` 
                        }}
                      >
                        <p className="font-semibold">Extracted Gradient Background</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Default Brand Colors Reference */}
        <div className="max-w-4xl mx-auto mt-12">
          <div className="bg-white rounded-2xl shadow-lg p-8">
            <h2 className="text-2xl font-bold mb-6 text-gray-800">Default Social Pay Brand Colors</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {Object.entries(SOCIAL_PAY_BRAND_COLORS).map(([name, color]) => (
                <div key={name} className="text-center">
                  <div 
                    className="w-full h-24 rounded-xl shadow-lg mb-4"
                    style={{ backgroundColor: color.hex }}
                  ></div>
                  <h3 className="font-semibold text-gray-800 mb-2 capitalize">{name}</h3>
                  <div className="space-y-1 text-sm text-gray-600">
                    <p className="font-mono bg-gray-100 px-2 py-1 rounded">{color.hex}</p>
                    <p>RGB({color.rgb.join(', ')})</p>
                    <p>HSL({color.hsl.join(', ')})</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 