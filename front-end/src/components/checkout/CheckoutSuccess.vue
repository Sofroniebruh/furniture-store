<template>
  <Wrapper class="py-8">
    <div class="max-w-2xl mx-auto text-center">
      <!-- Loading state -->
      <div v-if="isLoading" class="py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
        <p class="mt-4">Loading order details...</p>
      </div>
      
      <!-- Success state -->
      <div v-else-if="order" class="space-y-8">
        <!-- Success message -->
        <div class="text-center">
          <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <CheckCircle class="w-8 h-8 text-green-600" />
          </div>
          <h1 class="text-3xl font-bold text-green-600 mb-2">Payment Successful!</h1>
          <p class="text-gray-600">Thank you for your purchase. Your order has been confirmed.</p>
        </div>
        
        <!-- Order details -->
        <div class="bg-gray-50 border border-gray-200 rounded-lg p-6 text-left">
          <h2 class="text-xl font-semibold mb-4">Order Details</h2>
          
          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-600">Order ID:</span>
              <span class="font-mono text-sm">{{ order.id }}</span>
            </div>
            
            <div class="flex justify-between">
              <span class="text-gray-600">Order Date:</span>
              <span>{{ formatDate(order.created_at) }}</span>
            </div>
            
            <div class="flex justify-between">
              <span class="text-gray-600">Status:</span>
              <span class="capitalize px-2 py-1 bg-green-100 text-green-800 rounded text-sm">
                {{ order.status }}
              </span>
            </div>
            
            <div class="flex justify-between font-semibold text-lg">
              <span>Total Amount:</span>
              <span>{{ formatCurrency(order.total_amount) }}</span>
            </div>
          </div>
          
          <!-- Order items -->
          <div v-if="order.items && order.items.length > 0" class="mt-6">
            <h3 class="font-medium mb-3">Items Purchased:</h3>
            <div class="space-y-3">
              <div
                v-for="item in order.items"
                :key="item.id"
                class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded"
              >
                <div>
                  <h4 class="font-medium">{{ item.product?.name }}</h4>
                  <p class="text-sm text-gray-600">
                    {{ formatCurrency(item.price) }} × {{ item.quantity }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="font-medium">{{ formatCurrency(item.price * item.quantity) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Next steps -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 class="text-lg font-semibold text-blue-800 mb-2">What's Next?</h3>
          <ul class="text-blue-700 space-y-1 text-left">
            <li>• A confirmation email has been sent to your email address</li>
            <li>• Your items will be prepared for shipping</li>
            <li>• You'll receive tracking information once your order ships</li>
            <li>• Your items have been added to your purchase history</li>
          </ul>
        </div>
        
        <!-- Action buttons -->
        <div class="space-y-4">
          <router-link
            to="/dashboard/history"
            class="inline-block bg-[#c9a275] hover:bg-[#b8956a] text-white px-6 py-3 rounded-lg transition-colors"
          >
            View Purchase History
          </router-link>
          
          <div>
            <router-link
              to="/products"
              class="inline-block text-[#c9a275] hover:text-[#b8956a] px-6 py-2"
            >
              Continue Shopping
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- Error state -->
      <div v-else class="py-16">
        <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <X class="w-8 h-8 text-red-600" />
        </div>
        <h1 class="text-3xl font-bold text-red-600 mb-2">Order Not Found</h1>
        <p class="text-gray-600 mb-6">We couldn't find the order you're looking for.</p>
        <router-link
          to="/products"
          class="inline-block bg-[#c9a275] hover:bg-[#b8956a] text-white px-6 py-3 rounded-lg"
        >
          Continue Shopping
        </router-link>
      </div>
    </div>
  </Wrapper>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { CheckCircle, X } from 'lucide-vue-next'
import Wrapper from '@/components/Wrapper.vue'
import { useCartStore, type Order } from '@/stores/useCart'
import { useCart } from '@/composables/useCart'

const route = useRoute()
const cartStore = useCartStore()
const { loadCart } = useCart()

const isLoading = ref(true)
const order = ref<Order | null>(null)

const orderId = computed(() => route.query.order_id as string)

const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(amount)
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const loadOrderDetails = async () => {
  if (!orderId.value) {
    return
  }

  try {
    const orderData = await cartStore.getOrder(orderId.value)
    if (orderData) {
      order.value = orderData
      // Refresh cart (should be empty after successful payment)
      await loadCart()
    }
  } catch (err) {
    console.error('Error loading order:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadOrderDetails()
})
</script>