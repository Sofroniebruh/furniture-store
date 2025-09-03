<script setup lang="ts">
import { onMounted } from 'vue'
import { ShoppingCart, X } from 'lucide-vue-next'
import { useCart } from '@/composables/useCart'
import { useRouter } from 'vue-router'

const {
  cartItems,
  cartItemsCount,
  cartTotal,
  isLoading,
  error,
  formatCurrency,
  updateQuantity,
  removeItem,
  loadCart,
  proceedToCheckout
} = useCart()

const router = useRouter()

const emit = defineEmits<{
  close: []
}>()

const handleCheckout = async () => {
  const result = await proceedToCheckout()
  if (result.success && 'data' in result) {
    emit('close')
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
  console.log('CartDropdown: Component mounted, loading cart')
  loadCart()
})
</script>

<template>
  <div class="w-full bg-white">
    <div class="p-4">
      <h3 class="text-lg font-semibold mb-3">Shopping Cart</h3>
      
      <div v-if="isLoading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        <p class="text-sm text-gray-600 mt-2">Loading cart...</p>
      </div>
      
      <div v-else-if="error" class="text-center py-8 text-red-500">
        <X class="w-12 h-12 mx-auto mb-2 text-red-400" />
        <p class="text-sm mb-2">Failed to load cart</p>
        <p class="text-xs text-gray-600 mb-3">{{ error }}</p>
        <button
          @click="loadCart"
          class="bg-[#c9a275] hover:bg-[#b8956a] text-white px-3 py-1 rounded text-sm transition-colors"
        >
          Retry
        </button>
      </div>
      
      <div v-else-if="cartItemsCount === 0" class="text-center py-8 text-gray-500">
        <ShoppingCart class="w-12 h-12 mx-auto mb-2 text-gray-400" />
        <p>Your cart is empty</p>
      </div>
      
      <div v-else>
        <div class="space-y-3 max-h-64 overflow-y-auto">
          <div
            v-for="item in cartItems"
            :key="item.id"
            class="flex items-center space-x-3 p-2 hover:bg-gray-50 rounded"
          >
            <div class="flex-1">
              <h4 class="text-sm font-medium">{{ item.product?.name }}</h4>
              <p class="text-xs text-gray-600">{{ formatCurrency(item.product?.price || 0) }}</p>
              <div class="flex items-center mt-1">
                <button
                  @click="updateQuantity(item.id, item.quantity - 1)"
                  class="w-6 h-6 flex items-center justify-center bg-gray-100 hover:bg-gray-200 rounded text-sm"
                  :disabled="item.quantity <= 1"
                >
                  -
                </button>
                <span class="mx-2 text-sm">{{ item.quantity }}</span>
                <button
                  @click="updateQuantity(item.id, item.quantity + 1)"
                  class="w-6 h-6 flex items-center justify-center bg-gray-100 hover:bg-gray-200 rounded text-sm"
                >
                  +
                </button>
              </div>
            </div>
            <button
              @click="removeItem(item.id)"
              class="text-red-500 hover:text-red-700 p-1"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
        </div>
        
        <div class="border-t pt-3 mt-3">
          <div class="flex justify-between items-center mb-3">
            <span class="font-semibold">Total:</span>
            <span class="font-semibold">{{ formatCurrency(cartTotal) }}</span>
          </div>
          
          <div class="space-y-2">
            <router-link
              to="/cart"
              class="w-full bg-gray-100 hover:bg-gray-200 text-gray-800 py-2 px-4 rounded-md text-center block transition-colors"
              @click="$emit('close')"
            >
              View Cart
            </router-link>
            <button
              @click="handleCheckout"
              class="w-full bg-[#c9a275] hover:bg-[#b8956a] text-white py-2 px-4 rounded-md transition-colors cursor-pointer"
            >
              Checkout
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>