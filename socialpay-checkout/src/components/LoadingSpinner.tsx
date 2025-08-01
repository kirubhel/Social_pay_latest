import { ClipLoader } from 'react-spinners';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  color?: string;
  className?: string;
}

export default function LoadingSpinner({ 
  size = 'md', 
  color = '#30BB54',
  className = '' 
}: LoadingSpinnerProps) {
  const sizeMap = {
    sm: 16,
    md: 32, 
    lg: 48
  };

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <ClipLoader
        color={color}
        size={sizeMap[size]}
        aria-label="Loading Spinner"
        data-testid="loader"
      />
    </div>
  );
} 