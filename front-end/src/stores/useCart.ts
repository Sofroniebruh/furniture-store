import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CartProduct {
  id: string
  name: string
  description: string
  stock: number
  price: number
  pictureUrls: string[]
  event: string
  model: string
}

export interface CartItem {
  id: string
  user_id: string
  product_id: string
  quantity: number
  created_at: string
  updated_at: string
  product?: CartProduct
}

export interface Cart {
  user_id: string
  items: CartItem[]
  total_items: number
  total_price: number
}

export interface Order {
  id: string
  user_id: string
  stripe_payment_id: string
  total_amount: number
  status: 'pending' | 'paid' | 'cancelled' | 'refunded'
  created_at: string
  updated_at: string
  items?: OrderItem[]
  stripe_client_secret?: string
}

export interface OrderItem {
  id: string
  order_id: string
  product_id: string
  quantity: number
  price: number
  product?: CartProduct
}

export const useCartStore = defineStore('cart', () => {
  const cart = ref<Cart | null>(null)
  const isLoading = ref(false)
  const error = ref<string>('')

  const cartItemsCount = computed(() => cart.value?.total_items || 0)
  const cartTotal = computed(() => cart.value?.total_price || 0)
  const cartItems = computed(() => cart.value?.items || [])

  const fetchCart = async () => {
    isLoading.value = true
    error.value = ''

    try {
      console.log('Fetching cart from:', `${import.meta.env.VITE_CART_SERVICE_URL}/cart`)
      
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/cart`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      console.log('Cart fetch response status:', response.status)

      if (!response.ok) {
        const errorText = await response.text()
        console.error('Cart fetch error response:', errorText)
        throw new Error(`Failed to fetch cart: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('Cart data received:', data)
      cart.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch cart'
      console.error('Error fetching cart:', err)
      
      cart.value = {
        user_id: '',
        items: [],
        total_items: 0,
        total_price: 0
      }
    } finally {
      isLoading.value = false
    }
  }

  const addToCart = async (productId: string, quantity: number = 1) => {
    error.value = ''

    try {
      console.log('Adding to cart:', { product_id: productId, quantity })
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/cart/items`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          product_id: productId,
          quantity: quantity,
        }),
      })

      console.log('Add to cart response status:', response.status)

      if (!response.ok) {
        const errorData = await response.json()
        console.error('Add to cart error response:', errorData)
        throw new Error(errorData.error || 'Failed to add to cart')
      }

      await fetchCart()
      return { success: true }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to add to cart'
      console.error('Error adding to cart:', err)
      return { success: false, error: error.value }
    }
  }

  const updateCartItem = async (itemId: string, quantity: number) => {
    error.value = ''

    try {
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/cart/items/${itemId}`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          quantity: quantity,
        }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update cart item')
      }

      await fetchCart()
      return { success: true }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update cart item'
      console.error('Error updating cart item:', err)
      return { success: false, error: error.value }
    }
  }

  const removeFromCart = async (itemId: string) => {
    error.value = ''

    try {
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/cart/items/${itemId}`, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to remove from cart')
      }

      await fetchCart()
      return { success: true }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to remove from cart'
      console.error('Error removing from cart:', err)
      return { success: false, error: error.value }
    }
  }

  const clearCart = async () => {
    error.value = ''

    try {
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/cart`, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to clear cart')
      }

      cart.value = null
      return { success: true }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to clear cart'
      console.error('Error clearing cart:', err)
      return { success: false, error: error.value }
    }
  }

  const createCheckout = async () => {
    error.value = ''

    try {
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/checkout`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({}),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to create checkout')
      }

      const data = await response.json()
      return { success: true, data }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create checkout'
      console.error('Error creating checkout:', err)
      return { success: false, error: error.value }
    }
  }

  const getOrder = async (orderId: string): Promise<Order | null> => {
    error.value = ''

    try {
      const response = await fetch(`${import.meta.env.VITE_CART_SERVICE_URL}/orders/${orderId}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to fetch order')
      }

      const data = await response.json()
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch order'
      console.error('Error fetching order:', err)
      return null
    }
  }

  return {
    cart,
    isLoading,
    error,
    cartItemsCount,
    cartTotal,
    cartItems,
    fetchCart,
    addToCart,
    updateCartItem,
    removeFromCart,
    clearCart,
    createCheckout,
    getOrder,
  }
})