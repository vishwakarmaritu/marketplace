import React from 'react';

export default function BuyerLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      {/* navbar */}
      <nav className="bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center sticky top-0 z-50">
        <h1 className="text-2xl font-black text-blue-600">autopia</h1>
        <div className="flex gap-4">
          <button className="text-gray-600 hover:text-black font-medium">search</button>
          <button className="text-gray-600 hover:text-black font-medium">cart</button>
        </div>
      </nav>
      
      {/* content */}
      <main className="flex-grow p-6">
        {children}
      </main>
      
      {/* footer */}
      <footer className="bg-gray-900 text-white text-center p-6 mt-auto">
        <p className="text-sm">© 2026 autopia</p>
      </footer>
    </div>
  );
}