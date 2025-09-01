<template>
  <Wrapper class="py-8">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-3xl font-bold mb-8">Shopping Cart</h1>
      
      <div v-if="isLoading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
      
      <div v-else-if="cartItemsCount === 0" class="text-center py-16">
        <ShoppingCart class="w-24 h-24 mx-auto mb-4 text-gray-400" />
        <h2 class="text-xl font-semibold mb-2">Your cart is empty</h2>
        <p class="text-gray-600 mb-6">Start shopping to add items to your cart</p>
        <router-link
          to="/products"
          class="inline-block bg-[#c9a275] hover:bg-[#b8956a] text-white px-6 py-3 rounded-lg transition-colors"
        >
          Continue Shopping
        </router-link>
      </div>
      
      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-2 space-y-4">
          <div
            v-for="item in cartItems"
            :key="item.id"
            class="bg-white border border-gray-200 rounded-lg p-6"
          >
            <div class="flex items-start space-x-4">
              <div class="w-20 h-20 bg-gray-100 rounded-lg flex items-center justify-center">
                <div class="text-gray-400 text-sm">No image</div>
              </div>
              
              <div class="flex-1">
                <h3 class="text-lg font-semibold mb-1">{{ item.product?.name }}</h3>
                <p class="text-gray-600 text-sm mb-2">{{ item.product?.description }}</p>
                <p class="text-lg font-semibold text-[#c9a275]">
                  {{ formatCurrency(item.product?.price || 0) }}
                </p>
                
                <p v-if="item.product?.stock" class="text-sm text-gray-500 mt-1">
                  {{ item.product.stock }} in stock
                </p>
              </div>
              
              <div class="flex flex-col items-end space-y-2">
                <button
                  @click="removeItem(item.id)"
                  class="text-red-500 hover:text-red-700 p-1"
                  title="Remove item"
                >
                  <X class="w-5 h-5" />
                </button>
                
                <div class="flex items-center border border-gray-300 rounded">
                  <button
                    @click="updateQuantity(item.id, item.quantity - 1)"
                    class="px-3 py-1 hover:bg-gray-100"
                    :disabled="item.quantity <= 1"
                  >
                    -
                  </button>
                  <span class="px-3 py-1 bg-gray-50 min-w-[3rem] text-center">
                    {{ item.quantity }}
                  </span>
                  <button
                    @click="updateQuantity(item.id, item.quantity + 1)"
                    class="px-3 py-1 hover:bg-gray-100"
                    :disabled="item.quantity >= (item.product?.stock || 0)"
                  >
                    +
                  </button>
                </div>
                
                <p class="text-sm font-medium">
                  Subtotal: {{ formatCurrency(getCartItemTotal(item)) }}
                </p>
              </div>
            </div>
          </div>
          
          <div class="pt-4">
            <button
              @click="handleClearCart"
              class="text-red-600 hover:text-red-800 font-medium"
            >
              Clear Cart
            </button>
          </div>
        </div>
        
        <div class="lg:col-span-1">
          <div class="bg-gray-50 border border-gray-200 rounded-lg p-6 sticky top-24">
            <h3 class="text-lg font-semibold mb-4">Order Summary</h3>
            
            <div class="space-y-3">
              <div class="flex justify-between">
                <span>Items ({{ cartItemsCount }})</span>
                <span>{{ formatCurrency(cartTotal) }}</span>
              </div>
              
              <div class="flex justify-between">
                <span>Shipping</span>
                <span class="text-green-600">Free</span>
              </div>
              
              <div class="border-t pt-3">
                <div class="flex justify-between items-center text-lg font-semibold">
                  <span>Total</span>
                  <span>{{ formatCurrency(cartTotal) }}</span>
                </div>
              </div>
            </div>
            
            <button
              @click="handleCheckout"
              :disabled="cartItemsCount === 0"
              class="w-full mt-6 bg-[#c9a275] hover:bg-[#b8956a] disabled:bg-gray-300 disabled:cursor-not-allowed text-white py-3 px-4 rounded-lg font-medium transition-colors"
            >
              Proceed to Checkout
            </button>
            
            <router-link
              to="/products"
              class="block w-full mt-3 text-center text-[#c9a275] hover:text-[#b8956a] py-2"
            >
              Continue Shopping
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </Wrapper>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ShoppingCart, X } from 'lucide-vue-next'
import Wrapper from '@/components/Wrapper.vue'
import { useCart } from '@/composables/useCart'

const { 
  cartItems, 
  cartItemsCount, 
  cartTotal, 
  isLoading,
  formatCurrency, 
  updateQuantity, 
  removeItem,
  clearCart,
  loadCart,
  proceedToCheckout,
  getCartItemTotal
} = useCart()

const router = useRouter()

const handleClearCart = async () => {
  if (confirm('Are you sure you want to clear your cart?')) {
    await clearCart()
  }
}

const handleCheckout = async () => {
  const result = await proceedToCheckout()
  if (result.success) {
    router.push({
      name: 'checkout',
      query: {
        client_secret: result.data.client_secret,
        order_id: result.data.order_id,
      }
    })
  }
}

onMounted(() => {
  loadCart()
})
</script>