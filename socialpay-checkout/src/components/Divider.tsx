'use client'

interface DividerProps {
  className?: string;
}

export default function Divider({ className = '' }: DividerProps) {
  return (
    <div className={`w-full h-[1.5px] bg-gray-200 my-4 ${className}`} />
  );
} 