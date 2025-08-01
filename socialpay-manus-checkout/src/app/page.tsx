'use client';

import Image from "next/image";

export default function Home() {


  return (
    <div className="min-h-screen flex flex-col items-start justify-center p-8 max-w-7xl mx-auto">
      <div className="w-full max-w-2xl space-y-12">
        <div>
          <Image
            src="/socialpay.webp"
            alt="SocialPay logo"
            width={200}
            height={50}
            priority
          />
          <p className="mt-6 text-xl text-secondary/80">
            Connect. Accept. Grow.
          </p>
        </div>

        <div className="flex gap-4 flex-col sm:flex-row">
          <a
            href="https://onboarding.socialpay.co/"
            className="group relative flex items-center justify-center h-[52px] w-[204px] rounded-lg overflow-hidden"
          >
            {/* Animated gradient background */}
            <div className="absolute inset-0 bg-gradient-to-r from-green-400 via-emerald-500 to-teal-500 opacity-0 group-hover:opacity-100 transition-all duration-500 animate-gradient" />
            
            {/* Button content with backdrop blur */}
            <div className="relative flex items-center justify-center w-full h-full bg-primary group-hover:bg-transparent transition-colors duration-500 px-6 py-4">
              <span className="text-white font-medium relative z-10">Get Started</span>
            </div>
          </a>
          
          <a
            href="https://docs.socialpay.co"
            className="flex items-center justify-center h-[52px] w-[204px] rounded-lg bg-muted text-secondary font-medium hover:bg-muted/80 transition-all px-6 py-4"
          >
            Read Docs
          </a>
        </div>
      </div>
    </div>
  );
}
