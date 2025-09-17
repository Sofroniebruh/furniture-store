<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import type {Product} from '@/lib/types'
import DialogGeneral from "@/components/DialogGeneral.vue";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore";
import CreateDialogContent from "./CreateDialogContent.vue";
import EditProductDialog from "./EditProductDialog.vue";
import { useDashboard } from '@/composables/useDashboard'
import DeleteProductDialog from "@/components/dashboard-related/DeleteProductDialog.vue";

const { products, dashboardStats: apiDashboardStats, fetchAllProducts, deleteProduct: deleteProductApi, addProduct } = useDashboard()
const selectedProduct = ref<Product | null>(null)
const searchQuery = ref('')
const sortBy = ref('name')
const sortOrder = ref<'asc' | 'desc'>('asc')
const currentPage = ref(1)
const itemsPerPage = ref(10)

const stockHistory = ref<Array<{
  id: string
  product_id: string
  productName?: string
  type: 'in' | 'out' | 'adjustment'
  quantity: number
  previous_stock: number
  new_stock: number
  reason: string
  created_at: string
  product?: {
    name: string
    id: string
  }
}>>([])

onMounted(async () => {
  await fetchAllProducts()
  await loadStockHistory()
})

const filteredProducts = computed(() => {
  let filtered = products.value.filter(product =>
      product.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      product.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      product.model.toLowerCase().includes(searchQuery.value.toLowerCase())
  )

  filtered.sort((a, b) => {
    let aValue = a[sortBy.value as keyof Product]
    let bValue = b[sortBy.value as keyof Product]

    if (typeof aValue === 'string' && typeof bValue === 'string') {
      aValue = aValue.toLowerCase()
      bValue = bValue.toLowerCase()
    }
    else if (typeof aValue === 'number' && typeof bValue === 'number') {
      if (sortOrder.value === 'asc') {
        return aValue - bValue
      } else {
        return bValue - aValue
      }
    }

    if (sortOrder.value === 'asc') {
      return aValue > bValue ? 1 : -1
    } else {
      return aValue < bValue ? 1 : -1
    }
  })

  return filtered
})

const paginatedProducts = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredProducts.value.slice(start, end)
})

const totalPages = computed(() =>
    Math.ceil(filteredProducts.value.length / itemsPerPage.value)
)

const dashboardStats = computed(() => apiDashboardStats.value)

const openEditModal = (product: Product) => {
  selectedProduct.value = product
  setOpenDialog(true, {name: 'EditDialog'})
}

const openDeleteModal = (product: Product) => {
  selectedProduct.value = product
  setOpenDialog(true, {name: 'DeleteDialog'})
}

const closeModals = () => {
  selectedProduct.value = null
  setOpenDialog(false, {name: 'EditDialog'})
  setOpenDialog(false, {name: 'DeleteDialog'})
}

const deleteProduct = async () => {
  if (!selectedProduct.value) return

  const result = await deleteProductApi(selectedProduct.value.id)
  if (result.success) {
    closeModals()
  } else {
    console.error('Failed to delete product:', result.error)
  }
}

const handleProductAdded = (newProduct: Product) => {
  addProduct(newProduct)
  loadStockHistory()
}

const handleProductUpdated = (updatedProduct: Product) => {
  const index = products.value.findIndex(p => p.id === updatedProduct.id)
  if (index !== -1) {
    products.value[index] = updatedProduct
  }
  closeModals()
}

const updateStock = async (productId: string, quantity: number, reason: string) => {
  try {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/stock`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({
        product_id: productId,
        quantity: quantity,
        reason: reason
      })
    })
    
    if (!response.ok) {
      throw new Error(`Failed to update stock: ${response.status}`)
    }
    
    const result = await response.json()
    
    const productIndex = products.value.findIndex(p => p.id === productId)
    if (productIndex !== -1 && result.product) {
      products.value[productIndex] = result.product
    }
    
    await loadStockHistory()
    
  } catch (error) {
    console.error('Failed to update stock:', error)
  }
}

const loadStockHistory = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/stock-history?limit=20&page=1`, {
      method: 'GET',
      credentials: 'include'
    })
    
    if (!response.ok) {
      throw new Error(`Failed to load stock history: ${response.status}`)
    }
    
    const data = await response.json()
    
    stockHistory.value = data.stock_history || []
    
  } catch (error) {
    console.error('Failed to load stock history:', error)
    stockHistory.value = []
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount)
}

const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const getStockStatus = (amount: number) => {
  if (amount === 0) return {status: 'Out of Stock', class: 'text-red-600 bg-red-50'}
  if (amount < 10) return {status: 'Low Stock', class: 'text-orange-600 bg-orange-50'}
  return {status: 'In Stock', class: 'text-green-600 bg-green-50'}
}

const {setOpenDialog, isDialogOpen} = useScreenSheetStore()
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white shadow-sm border-b">
      <div class="px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div>
            <h1 class="text-xl sm:text-3xl font-bold text-gray-900">Dashboard</h1>
            <p class="text-sm sm:text-base text-gray-600 hidden sm:block">Manage your furniture inventory</p>
          </div>
          <div>
            <DialogGeneral class-name="h-[80%] flex items-center justify-center flex-col" title="Add new item" :is-open="isDialogOpen('CreationDialog')">
              <template #trigger>
                <button
                    @click="setOpenDialog(true, {name: 'CreationDialog'})"
                    class="bg-[#c9a275] hover:bg-[#dbb384] text-white px-3 py-2 sm:px-6 rounded-lg font-medium transition-colors text-sm sm:text-base"
                >
                  <span class="hidden sm:inline cursor-pointer">Add Product</span>
                  <span class="sm:hidden">Add</span>
                </button>
              </template>
              <template #content>
                <CreateDialogContent @product-added="handleProductAdded" />
              </template>
            </DialogGeneral>
          </div>
        </div>
      </div>
    </div>
    <div class="px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-6 mb-4 sm:mb-8">
        <div class="bg-white rounded-lg shadow p-3 sm:p-6">
          <div class="flex items-center">
            <div class="p-2 bg-blue-100 rounded-lg">
              <svg class="w-4 h-4 sm:w-6 sm:h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
              </svg>
            </div>
            <div class="ml-2 sm:ml-4">
              <p class="text-xs sm:text-sm font-medium text-gray-600">Products</p>
              <p class="text-lg sm:text-2xl font-bold text-gray-900">{{ dashboardStats.totalProducts }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-3 sm:p-6">
          <div class="flex items-center">
            <div class="p-2 bg-green-100 rounded-lg">
              <svg class="w-4 h-4 sm:w-6 sm:h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-14 0h14"></path>
              </svg>
            </div>
            <div class="ml-2 sm:ml-4">
              <p class="text-xs sm:text-sm font-medium text-gray-600">Stock</p>
              <p class="text-lg sm:text-2xl font-bold text-gray-900">{{ dashboardStats.totalStock }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-3 sm:p-6">
          <div class="flex items-center">
            <div class="p-2 bg-orange-100 rounded-lg">
              <svg class="w-4 h-4 sm:w-6 sm:h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
              </svg>
            </div>
            <div class="ml-2 sm:ml-4">
              <p class="text-xs sm:text-sm font-medium text-gray-600">Low Stock</p>
              <p class="text-lg sm:text-2xl font-bold text-gray-900">{{ dashboardStats.lowStock }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-3 sm:p-6">
          <div class="flex items-center">
            <div class="p-2 bg-purple-100 rounded-lg">
              <svg class="w-4 h-4 sm:w-6 sm:h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <div class="ml-2 sm:ml-4">
              <p class="text-xs sm:text-sm font-medium text-gray-600">Value</p>
              <p class="text-sm sm:text-2xl font-bold text-gray-900">{{ formatCurrency(dashboardStats.totalValue) }}</p>
            </div>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow mb-4 sm:mb-6">
        <div class="p-4 sm:p-6 border-b border-gray-200">
          <div class="space-y-3 sm:space-y-4">
            <div>
              <input
                  v-model="searchQuery"
                  type="text"
                  placeholder="Search products..."
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm sm:text-base"
              >
            </div>
            <div class="flex gap-2">
              <select
                  v-model="sortBy"
                  style="-webkit-appearance: none;"
                  class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              >
                <option value="name">Name</option>
                <option value="price">Price</option>
                <option value="stock">Stock</option>
                <option value="model">Model</option>
              </select>
              <button
                  @click="sortOrder = sortOrder === 'asc' ? 'desc' : 'asc'"
                  class="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path v-if="sortOrder === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M5 15l7-7 7 7"></path>
                  <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 9l-7 7-7-7"></path>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Model</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="product in paginatedProducts" :key="product.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <img
                        :src="product.pictureUrls[0] || '/images/furniture1.webp'"
                        :alt="product.name"
                        class="w-12 h-12 rounded-lg object-cover"
                    >
                    <div class="ml-4">
                      <div class="text-sm font-medium text-gray-900">{{ product.name }}</div>
                      <div class="text-sm text-gray-500">{{ product.description.substring(0, 50) }}...</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ product.model }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ formatCurrency(product.price) }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center space-x-2">
                    <span class="text-sm text-gray-900">{{ product.stock }}</span>
                    <button
                        @click="() => updateStock(product.id, 1, 'Manual stock increase')"
                        class="text-green-600 hover:text-green-800"
                        title="Add stock"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                              d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                      </svg>
                    </button>
                    <button
                        @click="() => updateStock(product.id, -1, 'Manual stock decrease')"
                        class="text-red-600 hover:text-red-800"
                        title="Remove stock"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"></path>
                      </svg>
                    </button>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <span :class="`px-2 py-1 text-xs font-medium rounded-full ${getStockStatus(product.stock).class}`">
                      {{ getStockStatus(product.stock).status }}
                    </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex space-x-2">
                    <DialogGeneral class-name="h-[80%] flex items-center justify-center flex-col" title="Edit Product" :is-open="isDialogOpen('EditDialog')">
                      <template #trigger>
                        <button
                            @click="openEditModal(product)"
                            class="text-blue-600 cursor-pointer hover:text-blue-900"
                        >
                          Edit
                        </button>
                      </template>
                      <template #content>
                        <EditProductDialog
                            :product="selectedProduct"
                            :is-open="true"
                            @close="closeModals"
                            @product-updated="handleProductUpdated">
                          </EditProductDialog>
                      </template>
                    </DialogGeneral>
                    <DialogGeneral class-name="flex items-center justify-center flex-col" title="Delete Product" :is-open="isDialogOpen('DeleteDialog')">
                      <template #trigger>
                        <button
                            @click="openDeleteModal(product)"
                            class="text-red-600 cursor-pointer hover:text-red-900"
                        >
                          Delete
                        </button>
                      </template>
                      <template #content>
                        <DeleteProductDialog
                            :product="selectedProduct"
                            @close="closeModals"
                            @delete-product="deleteProduct">
                        </DeleteProductDialog>
                      </template>
                    </DialogGeneral>
                  </div>
                </td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="px-4 sm:px-6 py-4 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <div class="text-xs sm:text-sm text-gray-700">
              <span class="hidden sm:inline">Showing {{
                  (currentPage - 1) * itemsPerPage + 1
                }} to {{ Math.min(currentPage * itemsPerPage, filteredProducts.length) }} of {{
                  filteredProducts.length
                }} results</span>
              <span class="sm:hidden">{{
                  (currentPage - 1) * itemsPerPage + 1
                }}-{{ Math.min(currentPage * itemsPerPage, filteredProducts.length) }} of {{
                  filteredProducts.length
                }}</span>
            </div>
            <div class="flex items-center space-x-2">
              <button
                  @click="currentPage = Math.max(1, currentPage - 1)"
                  :disabled="currentPage === 1"
                  class="px-2 sm:px-3 py-1 border border-gray-300 rounded-md text-xs sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <span class="hidden sm:inline">Previous</span>
                <span class="sm:hidden">‹</span>
              </button>
              <span class="px-2 sm:px-3 py-1 text-xs sm:text-sm text-gray-700">
                <span class="hidden sm:inline">Page {{ currentPage }} of {{ totalPages }}</span>
                <span class="sm:hidden">{{ currentPage }}/{{ totalPages }}</span>
              </span>
              <button
                  @click="currentPage = Math.min(totalPages, currentPage + 1)"
                  :disabled="currentPage === totalPages"
                  class="px-2 sm:px-3 py-1 border border-gray-300 rounded-md text-xs sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <span class="hidden sm:inline">Next</span>
                <span class="sm:hidden">›</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow">
        <div class="px-4 sm:px-6 py-4 border-b border-gray-200">
          <h3 class="text-base sm:text-lg font-medium text-gray-900">Recent Stock History</h3>
        </div>
        <div class="block sm:hidden">
          <div class="divide-y divide-gray-200">
            <div v-for="history in stockHistory.slice(0, 5)" :key="history.id" class="p-4">
              <div class="flex justify-between items-start mb-2">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">{{ history.product?.name || 'Unknown Product' }}</p>
                  <p class="text-xs text-gray-500">{{ formatDate(new Date(history.created_at)) }}</p>
                </div>
                <span :class="`px-2 py-1 text-xs font-medium rounded-full ${
                  history.type === 'in' ? 'text-green-600 bg-green-50' :
                  history.type === 'out' ? 'text-red-600 bg-red-50' :
                  'text-blue-600 bg-blue-50'
                }`">
                  {{ history.type.toUpperCase() }}
                </span>
              </div>
              <div class="flex justify-between items-center text-xs text-gray-600">
                <span>{{ history.quantity }} items</span>
                <span>{{ history.previous_stock }} → {{ history.new_stock }}</span>
              </div>
              <p class="text-xs text-gray-500 mt-1">{{ history.reason }}</p>
            </div>
          </div>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Quantity</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Previous</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">New</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Reason</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
            </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="history in stockHistory.slice(0, 10)" :key="history.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.product?.name || 'Unknown Product' }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="`px-2 py-1 text-xs font-medium rounded-full ${
                    history.type === 'in' ? 'text-green-600 bg-green-50' :
                    history.type === 'out' ? 'text-red-600 bg-red-50' :
                    'text-blue-600 bg-blue-50'
                  }`">
                    {{ history.type.toUpperCase() }}
                  </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.quantity }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.previous_stock }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.new_stock }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ history.reason }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(new Date(history.created_at)) }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>