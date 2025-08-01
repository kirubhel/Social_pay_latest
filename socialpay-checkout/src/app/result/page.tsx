'use client'

import { Suspense } from 'react'
import PaymentResultClient from './paymentresult'


export default function Page() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <PaymentResultClient/>
    </Suspense>
  )
}
