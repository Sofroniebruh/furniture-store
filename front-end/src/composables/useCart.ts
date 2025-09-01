import { computed } from 'vue'
import { useCartStore } from '@/stores/useCart'
import { useAuthStore } from '@/stores/useAuth'
import { toast } from 'vue-sonner'

export const useCart = () => {
  const cartStore = useCartStore()
  const authStore = useAuthStore()

  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const cartItemsCount = computed(() => cartStore.cartItemsCount)
  const cartTotal = computed(() => cartStore.cartTotal)
  const cartItems = computed(() => cartStore.cartItems)
  const isLoading = computed(() => cartStore.isLoading)
  const error = computed(() => cartStore.error)

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount)
  }

  const addToCart = async (productId: string, quantity: number = 1) => {
    if (!isAuthenticated.value) {
      toast('Please log in to add items to cart')
      return { success: false, error: 'Not authenticated' }
    }

    const result = await cartStore.addToCart(productId, quantity)
    
    if (result.success) {
      toast('Item added to cart')
    } else {
      toast(`Failed to add to cart: ${result.error}`)
    }

    return result
  }

  const updateQuantity = async (itemId: string, quantity: number) => {
    if (quantity <= 0) {
      return removeItem(itemId)
    }

    const result = await cartStore.updateCartItem(itemId, quantity)
    
    if (result.success) {
      toast('Cart updated')
    } else {
      toast(`Failed to update cart: ${result.error}`)
    }

    return result
  }

  const removeItem = async (itemId: string) => {
    const result = await cartStore.removeFromCart(itemId)
    
    if (result.success) {
      toast('Item removed from cart')
    } else {
      toast(`Failed to remove item: ${result.error}`)
    }

    return result
  }

  const clearCart = async () => {
    const result = await cartStore.clearCart()
    
    if (result.success) {
      toast('Cart cleared')
    } else {
      toast(`Failed to clear cart: ${result.error}`)
    }

    return result
  }

  const proceedToCheckout = async () => {
    if (!isAuthenticated.value) {
      toast('Please log in to checkout')
      return { success: false, error: 'Not authenticated' }
    }

    if (cartItemsCount.value === 0) {
      toast('Your cart is empty')
      return { success: false, error: 'Empty cart' }
    }

    const result = await cartStore.createCheckout()
    
    if (!result.success) {
      toast(`Checkout failed: ${result.error}`)
    }

    return result
  }

  const loadCart = async () => {
    console.log('useCart: loadCart called, isAuthenticated:', isAuthenticated.value)
    if (isAuthenticated.value) {
      console.log('useCart: User is authenticated, fetching cart')
      await cartStore.fetchCart()
    } else {
      console.log('useCart: User not authenticated, skipping cart fetch')
    }
  }

  const getItemQuantity = (productId: string): number => {
    const item = cartItems.value.find(item => item.product_id === productId)
    return item?.quantity || 0
  }

  const isItemInCart = (productId: string): boolean => {
    return cartItems.value.some(item => item.product_id === productId)
  }

  const getCartItemTotal = (item: any): number => {
    return (item.product?.price || 0) * item.quantity
  }

  return {
    // State
    cartItems,
    cartItemsCount,
    cartTotal,
    isLoading,
    error,
    isAuthenticated,

    // Actions
    addToCart,
    updateQuantity,
    removeItem,
    clearCart,
    loadCart,
    proceedToCheckout,

    // Utilities
    formatCurrency,
    getItemQuantity,
    isItemInCart,
    getCartItemTotal,

    // Store access
    cartStore,
  }
}