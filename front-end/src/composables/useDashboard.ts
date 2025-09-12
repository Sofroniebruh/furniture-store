import { ref, computed } from 'vue'
import type { Product } from '@/lib/types'

export function useDashboard() {
  const products = ref<Product[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  const dashboardStats = computed(() => ({
    totalProducts: products.value.length,
    totalStock: products.value.reduce((sum, product) => sum + product.stock, 0),
    lowStock: products.value.filter(product => product.stock < 10).length,
    outOfStock: products.value.filter(product => product.stock === 0).length,
    totalValue: products.value.reduce((sum, product) => sum + (product.price * product.stock), 0)
  }))

  const fetchAllProducts = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products?limit=1000&page=1`, {
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        throw new Error(`Failed to fetch products: ${response.status}`)
      }
      
      const data = await response.json()
      products.value = data.products || []
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch products'
      console.error('Dashboard fetch error:', err)
    } finally {
      loading.value = false
    }
  }

  const deleteProduct = async (productId: string) => {
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products?id=${productId}`, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        throw new Error(`Failed to delete product: ${response.status}`)
      }
      
      const index = products.value.findIndex(p => p.id === productId)
      if (index !== -1) {
        products.value.splice(index, 1)
      }
      
      return { success: true }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete product'
      console.error('Delete product error:', err)
      return { success: false, error: errorMessage }
    }
  }

  const updateProduct = async (productId: string, formData: FormData) => {
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products?id=${productId}`, {
        method: 'PUT',
        credentials: 'include',
        body: formData
      })
      
      if (!response.ok) {
        throw new Error(`Failed to update product: ${response.status}`)
      }
      
      const data = await response.json()
      
      const index = products.value.findIndex(p => p.id === productId)
      if (index !== -1) {
        products.value[index] = data.updated
      }
      
      return { success: true, product: data.updated }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update product'
      console.error('Update product error:', err)
      return { success: false, error: errorMessage }
    }
  }

  const addProduct = (newProduct: Product) => {
    products.value.push(newProduct)
  }

  return {
    products,
    loading,
    error,
    dashboardStats,
    fetchAllProducts,
    deleteProduct,
    updateProduct,
    addProduct
  }
}