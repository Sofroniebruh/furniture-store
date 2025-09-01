<template>
  <Wrapper class="py-8">
    <div class="max-w-2xl mx-auto">
      <h1 class="text-3xl font-bold mb-8">Checkout</h1>
      
      <!-- Loading state -->
      <div v-if="isLoading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
      
      <!-- Error state -->
      <div v-else-if="error" class="text-center py-16">
        <div class="text-red-600 mb-4">{{ error }}</div>
        <router-link
          to="/cart"
          class="inline-block bg-[#c9a275] hover:bg-[#b8956a] text-white px-6 py-3 rounded-lg"
        >
          Back to Cart
        </router-link>
      </div>
      
      <!-- Checkout form -->
      <div v-else class="space-y-8">
        <!-- Order summary -->
        <div class="bg-gray-50 border border-gray-200 rounded-lg p-6">
          <h2 class="text-xl font-semibold mb-4">Order Summary</h2>
          
          <div v-if="order" class="space-y-3">
            <div class="flex justify-between">
              <span>Total Amount</span>
              <span class="font-semibold">{{ formatCurrency(order.total_amount) }}</span>
            </div>
            
            <div v-if="order.items && order.items.length > 0">
              <h3 class="font-medium mb-2">Items:</h3>
              <div class="space-y-2">
                <div
                  v-for="item in order.items"
                  :key="item.id"
                  class="flex justify-between text-sm"
                >
                  <span>{{ item.product?.name }} × {{ item.quantity }}</span>
                  <span>{{ formatCurrency(item.price * item.quantity) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Payment form -->
        <div class="bg-white border border-gray-200 rounded-lg p-6">
          <h2 class="text-xl font-semibold mb-4">Payment Information</h2>
          
          <div v-if="!stripeLoaded" class="text-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
            <p class="mt-2">Loading payment form...</p>
          </div>
          
          <div v-else>
            <!-- Stripe Elements will be mounted here -->
            <div id="payment-element" class="mb-6"></div>
            
            <div v-if="paymentError" class="text-red-600 text-sm mb-4">
              {{ paymentError }}
            </div>
            
            <button
              @click="handleSubmit"
              :disabled="isProcessing || !paymentElement"
              class="w-full bg-[#c9a275] hover:bg-[#b8956a] disabled:bg-gray-300 disabled:cursor-not-allowed text-white py-3 px-4 rounded-lg font-medium transition-colors"
            >
              <span v-if="isProcessing">Processing...</span>
              <span v-else>Complete Payment</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Wrapper>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Wrapper from '@/components/Wrapper.vue'
import { useCartStore, type Order } from '@/stores/useCart'
import { toast } from 'vue-sonner'

// Stripe types
declare global {
  interface Window {
    Stripe: any
  }
}

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()

const isLoading = ref(true)
const error = ref('')
const order = ref<Order | null>(null)
const stripe = ref<any>(null)
const elements = ref<any>(null)
const paymentElement = ref<any>(null)
const stripeLoaded = ref(false)
const isProcessing = ref(false)
const paymentError = ref('')

const clientSecret = computed(() => route.query.client_secret as string)
const orderId = computed(() => route.query.order_id as string)

const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(amount)
}

const loadStripe = async () => {
  try {
    // Load Stripe.js
    if (!window.Stripe) {
      const script = document.createElement('script')
      script.src = 'https://js.stripe.com/v3/'
      script.onload = initializeStripe
      document.head.appendChild(script)
    } else {
      initializeStripe()
    }
  } catch (err) {
    error.value = 'Failed to load payment system'
    console.error('Error loading Stripe:', err)
  }
}

const initializeStripe = async () => {
  try {
    if (!clientSecret.value) {
      error.value = 'Invalid checkout session'
      return
    }

    // Initialize Stripe
    stripe.value = window.Stripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)
    
    if (!stripe.value) {
      throw new Error('Failed to initialize Stripe')
    }

    // Create Elements instance
    elements.value = stripe.value.elements({
      clientSecret: clientSecret.value,
      appearance: {
        theme: 'stripe',
        variables: {
          colorPrimary: '#c9a275',
        }
      }
    })

    // Create and mount Payment Element
    paymentElement.value = elements.value.create('payment')
    paymentElement.value.mount('#payment-element')

    paymentElement.value.on('ready', () => {
      stripeLoaded.value = true
    })

    paymentElement.value.on('change', (event: any) => {
      if (event.error) {
        paymentError.value = event.error.message
      } else {
        paymentError.value = ''
      }
    })
  } catch (err) {
    error.value = 'Failed to initialize payment form'
    console.error('Error initializing Stripe:', err)
  }
}

const handleSubmit = async () => {
  if (!stripe.value || !elements.value) {
    return
  }

  isProcessing.value = true
  paymentError.value = ''

  try {
    const { error: submitError } = await elements.value.submit()
    if (submitError) {
      paymentError.value = submitError.message
      return
    }

    const { error: confirmError, paymentIntent } = await stripe.value.confirmPayment({
      elements: elements.value,
      confirmParams: {
        return_url: `${window.location.origin}/checkout/success?order_id=${orderId.value}`,
      },
      redirect: 'if_required'
    })

    if (confirmError) {
      paymentError.value = confirmError.message
    } else if (paymentIntent && paymentIntent.status === 'succeeded') {
      // Payment succeeded, redirect to success page
      toast('Payment successful!')
      router.push({
        name: 'checkout-success',
        query: { order_id: orderId.value }
      })
    }
  } catch (err) {
    paymentError.value = 'An unexpected error occurred'
    console.error('Error processing payment:', err)
  } finally {
    isProcessing.value = false
  }
}

const loadOrderDetails = async () => {
  if (!orderId.value) {
    error.value = 'Invalid order ID'
    return
  }

  try {
    const orderData = await cartStore.getOrder(orderId.value)
    if (orderData) {
      order.value = orderData
    } else {
      error.value = 'Order not found'
    }
  } catch (err) {
    error.value = 'Failed to load order details'
    console.error('Error loading order:', err)
  }
}

onMounted(async () => {
  try {
    await loadOrderDetails()
    await loadStripe()
  } catch (err) {
    error.value = 'Failed to initialize checkout'
    console.error('Error in checkout initialization:', err)
  } finally {
    isLoading.value = false
  }
})
</script>