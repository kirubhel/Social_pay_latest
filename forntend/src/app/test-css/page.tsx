export default function TestCSSPage() {
  return (
    <div className="min-h-screen bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">CSS Test Page</h1>
        <p className="text-gray-600 mb-6">
          If you can see styled colors, gradients, and this card layout, 
          then Tailwind CSS is working correctly!
        </p>
        <div className="space-y-4">
          <button className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors">
            Blue Button
          </button>
          <button className="w-full bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition-colors">
            Green Button
          </button>
          <button className="w-full bg-orange-500 hover:bg-orange-600 text-white font-medium py-2 px-4 rounded-md transition-colors">
            Orange Button
          </button>
        </div>
        <div className="mt-6 p-4 bg-gray-100 rounded-md">
          <p className="text-sm text-gray-700">
            ✅ Background gradient<br/>
            ✅ White card with shadow<br/>
            ✅ Colored buttons with hover effects<br/>
            ✅ Responsive typography
          </p>
        </div>
      </div>
    </div>
  )
} 