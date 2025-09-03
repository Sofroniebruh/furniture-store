<script setup lang="ts">
import { onMounted } from 'vue'
import { ShoppingCart, X } from 'lucide-vue-next'
import { useCart } from '@/composables/useCart'
import { useRouter } from 'vue-router'

defineEmits<{
  close: []
}>()

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
  proceedToCheckout,
  getCartItemTotal
} = useCart()

const router = useRouter()

const handleCheckout = async () => {
  const result = await proceedToCheckout()
  if (result.success && 'data' in result) {
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

<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold">Shopping Cart</h3>
      <button @click="$emit('close')" class="p-1 hover:bg-gray-100 rounded">
        <X class="w-5 h-5" />
      </button>
    </div>
    
    <div v-if="isLoading" class="flex-1 flex justify-center items-center">
      <div class="text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-2"></div>
        <p class="text-sm text-gray-600">Loading cart...</p>
      </div>
    </div>
    
    <div v-else-if="error" class="flex-1 flex flex-col justify-center items-center text-red-500">
      <div class="text-center">
        <X class="w-16 h-16 mb-4 text-red-400 mx-auto" />
        <p class="text-lg font-medium mb-2">Failed to load cart</p>
        <p class="text-sm text-center mb-4">{{ error }}</p>
        <button
          @click="loadCart"
          class="bg-[#c9a275] hover:bg-[#b8956a] text-white px-4 py-2 rounded transition-colors"
        >
          Retry
        </button>
      </div>
    </div>
    
    <div v-else-if="cartItemsCount === 0" class="flex-1 flex flex-col justify-center items-center text-gray-500">
      <ShoppingCart class="w-16 h-16 mb-4 text-gray-400" />
      <p class="text-lg font-medium mb-2">Your cart is empty</p>
      <p class="text-sm text-center mb-4">Add some products to get started!</p>
      <router-link
        to="/products"
        @click="$emit('close')"
        class="bg-[#c9a275] hover:bg-[#b8956a] text-white px-4 py-2 rounded transition-colors"
      >
        Shop Now
      </router-link>
    </div>
    
    <div v-else class="flex-1 flex flex-col">
      <div class="flex-1 overflow-y-auto space-y-3 mb-4">
        <div
          v-for="item in cartItems"
          :key="item.id"
          class="flex items-center space-x-3 p-3"
        >
          <div class="w-12 h-12 bg-gray-100 rounded flex items-center justify-center flex-shrink-0">
            <div class="text-xs text-gray-400">IMG</div>
          </div>
          
          <div class="flex-1 min-w-0">
            <h4 class="text-sm font-medium truncate">{{ item.product?.name }}</h4>
            <p class="text-xs text-gray-600">{{ formatCurrency(item.product?.price || 0) }}</p>
            <div class="flex items-center mt-2">
              <button
                @click="updateQuantity(item.id, item.quantity - 1)"
                class="w-6 h-6 flex items-center justify-center bg-gray-100 hover:bg-gray-200 rounded text-sm"
                :disabled="item.quantity <= 1"
              >
                -
              </button>
              <span class="mx-2 text-sm min-w-[2rem] text-center">{{ item.quantity }}</span>
              <button
                @click="updateQuantity(item.id, item.quantity + 1)"
                class="w-6 h-6 flex items-center justify-center bg-gray-100 hover:bg-gray-200 rounded text-sm"
                :disabled="item.quantity >= (item.product?.stock || 0)"
              >
                +
              </button>
            </div>
          </div>
          
          <div class="flex flex-col items-end">
            <button
              @click="removeItem(item.id)"
              class="text-red-500 hover:text-red-700 p-1 mb-1"
            >
              <X class="w-4 h-4" />
            </button>
            <p class="text-sm font-medium">
              {{ formatCurrency(getCartItemTotal(item)) }}
            </p>
          </div>
        </div>
      </div>
      
      <div class="border-t pt-4 bg-white">
        <div class="flex justify-between items-center mb-4">
          <div>
            <span class="text-sm text-gray-600">{{ cartItemsCount }} items</span>
          </div>
          <span class="text-lg font-semibold">{{ formatCurrency(cartTotal) }}</span>
        </div>
        
        <div class="space-y-2">
          <router-link
            to="/cart"
            class="w-full bg-gray-100 hover:bg-gray-200 text-gray-800 py-3 px-4 rounded-md text-center block transition-colors font-medium"
            @click="$emit('close')"
          >
            View Cart
          </router-link>
          <button
            @click="handleCheckout"
            class="w-full bg-[#c9a275] hover:bg-[#b8956a] text-white py-3 px-4 rounded-md transition-colors font-medium"
            :disabled="cartItemsCount === 0"
          >
            Checkout
          </button>
        </div>
      </div>
    </div>
  </div>
</template>