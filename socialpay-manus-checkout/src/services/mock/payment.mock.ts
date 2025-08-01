import { PaymentResponse } from '../payment.service';

// Mock successful transaction response with redirect (for card payments)
const mockCardPaymentResponse: PaymentResponse = {
  success: true,
  data: {
    status: {
      current: "TXN_INIT",
      next: "TXN_PROCESS"
    },
    transaction: {
      id: "91EADE46-AB61-4DFC-B64F-E0990722F0A0",
      pricing: {
        amount: 100.00,
        fees: [
          {
            transaction: 5.00
          }
        ],
        total_amount: 105.00
      },
      status: {
        value: "pending",
        msg: "Transaction is initiated, waiting for user's confirmation"
      }
    }
  }
};

// Mock successful process response for card payments
const mockCardProcessResponse: PaymentResponse = {
  success: true,
  data: {
    status: {
      current: "TXN_PROCESS",
      next: "TXN_CHECKOUT"
    },
    checkout: {
      type: "redirect",
      data: "http://196.190.251.68:8008/checkout/vault/91EADE46-AB61-4DFC-B64F-E0990722F0A0.html"
    }
  }
};

// Mock successful transaction response for wallet payments (Telebirr, etc.)
const mockWalletPaymentResponse: PaymentResponse = {
  success: true,
  data: {
    status: {
      current: "TXN_INIT",
      next: "TXN_PROCESS"
    },
    transaction: {
      id: "B2C45F67-DE89-4ABC-9012-34567890ABCD",
      pricing: {
        amount: 100.00,
        fees: [
          {
            transaction: 2.50
          }
        ],
        total_amount: 102.50
      },
      status: {
        value: "pending",
        msg: "Transaction is initiated, waiting for wallet confirmation"
      }
    }
  }
};

// Mock successful process response for wallet payments
const mockWalletProcessResponse: PaymentResponse = {
  success: true,
  data: {
    status: {
      current: "TXN_PROCESS",
      next: "TXN_CHECKOUT"
    },
    checkout: {
      type: "notification",
      data: "Please check your phone for payment confirmation"
    }
  }
};

// Mock successful transaction response for bank payments
const mockBankPaymentResponse: PaymentResponse = {
  success: true,
  data: {
    status: {
      current: "TXN_INIT",
      next: "TXN_PROCESS"
    },
    transaction: {
      id: "C3D56F78-EF90-4BCD-A123-45678901BCDE",
      pricing: {
        amount: 100.00,
        fees: [
          {
            transaction: 3.00
          }
        ],
        total_amount: 103.00
      },
      status: {
        value: "pending",
        msg: "Transaction is initiated, waiting for bank confirmation"
      }
    }
  }
};

// Mock successful process response for bank payments
const mockBankProcessResponse: PaymentResponse = {
  success: true,
  data: {
    status: {
      current: "TXN_PROCESS",
      next: "TXN_CHECKOUT"
    },
    checkout: {
      type: "notification",
      data: "Please check your phone for bank confirmation"
    }
  }
};

// Mock error responses
const mockErrorResponses = {
  network: {
    success: false,
    error: {
      message: "Network error occurred. Please check your connection.",
      statusCode: 0,
      type: "NETWORK_ERROR"
    }
  },
  invalid: {
    success: false,
    error: {
      message: "Invalid payment details provided.",
      statusCode: 400,
      type: "BAD_REQUEST"
    }
  },
  server: {
    success: false,
    error: {
      message: "Server error occurred. Please try again later.",
      statusCode: 500,
      type: "USER_ERROR"
    }
  }
};

// Helper function to simulate API delay
const simulateDelay = () => new Promise(resolve => setTimeout(resolve, 1000));

// Mock payment service functions
export const mockPaymentService = {
  initiatePayment: async (medium: string) => {
    await simulateDelay();
    
    // Randomly throw errors (20% chance)
    if (Math.random() < 0.2) {
      const errors = Object.values(mockErrorResponses);
      throw errors[Math.floor(Math.random() * errors.length)];
    }

    switch (medium) {
      case 'CYBERSOURCE':
        return mockCardPaymentResponse;
      case 'TELEBIRR':
      case 'CBE':
      case 'MPESA':
        return mockWalletPaymentResponse;
      case 'AWINETAA':
      case 'BUNAETAA':
        return mockBankPaymentResponse;
      default:
        throw mockErrorResponses.invalid;
    }
  },

  processPayment: async (medium: string) => {
    await simulateDelay();
    
    // Randomly throw errors (10% chance)
    if (Math.random() < 0.1) {
      const errors = Object.values(mockErrorResponses);
      throw errors[Math.floor(Math.random() * errors.length)];
    }

    switch (medium) {
      case 'CYBERSOURCE':
        return mockCardProcessResponse;
      case 'TELEBIRR':
      case 'CBE':
      case 'MPESA':
        return mockWalletProcessResponse;
      case 'AWINETAA':
      case 'BUNAETAA':
        return mockBankProcessResponse;
      default:
        throw mockErrorResponses.invalid;
    }
  }
}; 