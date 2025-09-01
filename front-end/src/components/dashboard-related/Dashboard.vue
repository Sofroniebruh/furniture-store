<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {productsData} from '@/lib/data'
import type {Product} from '@/lib/types'
import DialogGeneral from "@/components/DialogGeneral.vue";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore";

const products = ref<Product[]>([])
const selectedProduct = ref<Product | null>(null)
const isEditModalOpen = ref(false)
const isAddModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const searchQuery = ref('')
const sortBy = ref('name')
const sortOrder = ref<'asc' | 'desc'>('asc')
const currentPage = ref(1)
const itemsPerPage = ref(10)
const isMobileMenuOpen = ref(false)
const viewMode = ref<'cards' | 'table'>('cards') // Default to cards for better mobile experience

// Form data for adding/editing products
const formData = reactive({
  name: '',
  description: '',
  price: 0,
  amount: 0,
  model: '',
  event: 'none',
  pictureUrls: [''],
  colors: [{id: '', name: ''}]
})

// Stock management
const stockHistory = ref<Array<{
  id: string
  productId: string
  productName: string
  type: 'in' | 'out' | 'adjustment'
  quantity: number
  previousStock: number
  newStock: number
  reason: string
  date: Date
}>>([])

// Initialize data
onMounted(() => {
  products.value = [...productsData]
  loadStockHistory()

  // Set view mode based on screen size
  if (window.innerWidth < 768) {
    viewMode.value = 'cards'
  }
})

// Computed properties
const filteredProducts = computed(() => {
  let filtered = products.value.filter(product =>
      product.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      product.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      product.model.toLowerCase().includes(searchQuery.value.toLowerCase())
  )

  // Sorting
  filtered.sort((a, b) => {
    let aValue = a[sortBy.value as keyof Product]
    let bValue = b[sortBy.value as keyof Product]

    if (typeof aValue === 'string' && typeof bValue === 'string') {
      aValue = aValue.toLowerCase()
      bValue = bValue.toLowerCase()
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

const dashboardStats = computed(() => ({
  totalProducts: products.value.length,
  totalStock: products.value.reduce((sum, product) => sum + product.amount, 0),
  lowStock: products.value.filter(product => product.amount < 10).length,
  outOfStock: products.value.filter(product => product.amount === 0).length,
  totalValue: products.value.reduce((sum, product) => sum + (product.price * product.amount), 0)
}))

// CRUD Operations
const openAddModal = () => {
  resetForm()
  isAddModalOpen.value = true
  document.body.style.overflow = 'hidden' // Prevent background scroll
}

const openEditModal = (product: Product) => {
  selectedProduct.value = product
  formData.name = product.name
  formData.description = product.description
  formData.price = product.price
  formData.amount = product.amount
  formData.model = product.model
  formData.event = product.event
  formData.pictureUrls = [...product.pictureUrls]
  formData.colors = [...product.colors]
  isEditModalOpen.value = true
  document.body.style.overflow = 'hidden'
}

const openDeleteModal = (product: Product) => {
  selectedProduct.value = product
  isDeleteModalOpen.value = true
  document.body.style.overflow = 'hidden'
}

const closeModals = () => {
  isAddModalOpen.value = false
  isEditModalOpen.value = false
  isDeleteModalOpen.value = false
  selectedProduct.value = null
  document.body.style.overflow = 'auto'
}

const resetForm = () => {
  formData.name = ''
  formData.description = ''
  formData.price = 0
  formData.amount = 0
  formData.model = ''
  formData.event = 'none'
  formData.pictureUrls = ['']
  formData.colors = [{id: '', name: ''}]
}

const addProduct = () => {
  const newProduct: Product = {
    id: generateId(),
    name: formData.name,
    description: formData.description,
    price: formData.price,
    amount: formData.amount,
    model: formData.model,
    event: formData.event,
    pictureUrls: formData.pictureUrls.filter(url => url.trim()),
    colors: formData.colors.filter(color => color.name.trim())
  }

  products.value.push(newProduct)
  addStockHistory(newProduct.id, newProduct.name, 'in', newProduct.amount, 0, newProduct.amount, 'Initial stock')
  closeModals()
  resetForm()
}

const updateProduct = () => {
  if (!selectedProduct.value) return

  const index = products.value.findIndex(p => p.id === selectedProduct.value!.id)
  if (index !== -1) {
    const oldAmount = products.value[index].amount
    const newAmount = formData.amount

    products.value[index] = {
      ...selectedProduct.value,
      name: formData.name,
      description: formData.description,
      price: formData.price,
      amount: newAmount,
      model: formData.model,
      event: formData.event,
      pictureUrls: formData.pictureUrls.filter(url => url.trim()),
      colors: formData.colors.filter(color => color.name.trim())
    }

    if (newAmount !== oldAmount) {
      const type = newAmount > oldAmount ? 'in' : 'out'
      const quantity = Math.abs(newAmount - oldAmount)
      addStockHistory(
          selectedProduct.value.id,
          formData.name,
          type,
          quantity,
          oldAmount,
          newAmount,
          'Stock adjustment'
      )
    }
  }

  closeModals()
}

const deleteProduct = () => {
  if (!selectedProduct.value) return

  const index = products.value.findIndex(p => p.id === selectedProduct.value!.id)
  if (index !== -1) {
    products.value.splice(index, 1)
  }

  closeModals()
}

// Stock Management
const updateStock = (productId: string, quantity: number, reason: string) => {
  const product = products.value.find(p => p.id === productId)
  if (!product) return

  const oldAmount = product.amount
  const newAmount = Math.max(0, oldAmount + quantity)
  const type = quantity > 0 ? 'in' : 'out'

  product.amount = newAmount

  addStockHistory(
      productId,
      product.name,
      type,
      Math.abs(quantity),
      oldAmount,
      newAmount,
      reason
  )
}

const addStockHistory = (
    productId: string,
    productName: string,
    type: 'in' | 'out' | 'adjustment',
    quantity: number,
    previousStock: number,
    newStock: number,
    reason: string
) => {
  stockHistory.value.unshift({
    id: generateId(),
    productId,
    productName,
    type,
    quantity,
    previousStock,
    newStock,
    reason,
    date: new Date()
  })
}

const loadStockHistory = () => {
  // Simulate loading stock history
  stockHistory.value = [
    {
      id: '1',
      productId: products.value[0]?.id || '',
      productName: products.value[0]?.name || '',
      type: 'in',
      quantity: 50,
      previousStock: 0,
      newStock: 50,
      reason: 'Initial stock',
      date: new Date('2024-01-01')
    }
  ]
}

// Utility functions
const generateId = () => {
  return Math.random().toString(36).substr(2, 9)
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

// Form helpers
const addColor = () => {
  formData.colors.push({id: generateId(), name: ''})
}

const removeColor = (index: number) => {
  formData.colors.splice(index, 1)
}

const addPictureUrl = () => {
  formData.pictureUrls.push('')
}

const removePictureUrl = (index: number) => {
  formData.pictureUrls.splice(index, 1)
}

//-------------
const {setOpenDialog, isDialogOpen} = useScreenSheetStore()
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div>
            <h1 class="text-xl sm:text-3xl font-bold text-gray-900">Dashboard</h1>
            <p class="text-sm sm:text-base text-gray-600 hidden sm:block">Manage your furniture inventory</p>
          </div>
          <div>
            <DialogGeneral :is-open="isDialogOpen('CreationDialog')">
              <template #trigger>
                <button
                    @click="setOpenDialog(true, {name: 'CreationDialog'})"
                    class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 sm:px-6 rounded-lg font-medium transition-colors text-sm sm:text-base"
                >
                  <span class="hidden sm:inline">Add Product</span>
                  <span class="sm:hidden">Add</span>
                </button>
              </template>
              <template #content>
                <p>damn</p>
              </template>
            </DialogGeneral>

          </div>
        </div>
      </div>
    </div>

    <!-- Stats Cards - Mobile Optimized -->
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

      <!-- Search and Filters - Mobile Optimized -->
      <div class="bg-white rounded-lg shadow mb-4 sm:mb-6">
        <div class="p-4 sm:p-6 border-b border-gray-200">
          <div class="space-y-3 sm:space-y-4">
            <!-- Search -->
            <div>
              <input
                  v-model="searchQuery"
                  type="text"
                  placeholder="Search products..."
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm sm:text-base"
              >
            </div>
            <!-- Filters -->
            <div class="flex gap-2">
              <select
                  v-model="sortBy"
                  class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              >
                <option value="name">Name</option>
                <option value="price">Price</option>
                <option value="amount">Stock</option>
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
                    <span class="text-sm text-gray-900">{{ product.amount }}</span>
                    <button
                        @click="updateStock(product.id, 1, 'Manual stock increase')"
                        class="text-green-600 hover:text-green-800"
                        title="Add stock"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                              d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                      </svg>
                    </button>
                    <button
                        @click="updateStock(product.id, -1, 'Manual stock decrease')"
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
                    <span :class="`px-2 py-1 text-xs font-medium rounded-full ${getStockStatus(product.amount).class}`">
                      {{ getStockStatus(product.status).status }}
                    </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex space-x-2">
                    <button
                        @click="openEditModal(product)"
                        class="text-blue-600 hover:text-blue-900"
                    >
                      Edit
                    </button>
                    <button
                        @click="openDeleteModal(product)"
                        class="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Mobile-Optimized Pagination -->
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

      <!-- Mobile-Optimized Stock History -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-4 sm:px-6 py-4 border-b border-gray-200">
          <h3 class="text-base sm:text-lg font-medium text-gray-900">Recent Stock History</h3>
        </div>

        <!-- Mobile Card View for Stock History -->
        <div class="block sm:hidden">
          <div class="divide-y divide-gray-200">
            <div v-for="history in stockHistory.slice(0, 5)" :key="history.id" class="p-4">
              <div class="flex justify-between items-start mb-2">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">{{ history.productName }}</p>
                  <p class="text-xs text-gray-500">{{ formatDate(history.date) }}</p>
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
                <span>{{ history.previousStock }} → {{ history.newStock }}</span>
              </div>
              <p class="text-xs text-gray-500 mt-1">{{ history.reason }}</p>
            </div>
          </div>
        </div>

        <!-- Desktop Table View for Stock History -->
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
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.productName }}</td>
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
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.previousStock }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ history.newStock }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ history.reason }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(history.date) }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Mobile-Optimized Add Product Modal -->
    <div v-if="isAddModalOpen"
         class="fixed inset-0 bg-black bg-opacity-50 flex items-end sm:items-center justify-center z-50"
         @click.self="closeModals">
      <div
          class="bg-white rounded-t-lg sm:rounded-lg shadow-xl w-full sm:max-w-2xl sm:mx-4 max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white px-4 sm:px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">Add New Product</h3>
          <button
              @click="closeModals"
              class="p-2 hover:bg-gray-100 rounded-lg"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        <form @submit.prevent="addProduct" class="p-4 sm:p-6 space-y-4 sm:space-y-6">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
            <div class="sm:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Name</label>
              <input
                  v-model="formData.name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Model</label>
              <input
                  v-model="formData.model"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Event</label>
              <select
                  v-model="formData.event"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="none">None</option>
                <option value="sale">Sale</option>
                <option value="featured">Featured</option>
                <option value="new">New</option>
                <option value="bestseller">Bestseller</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Price ($)</label>
              <input
                  v-model="formData.price"
                  type="number"
                  step="0.01"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Initial Stock</label>
              <input
                  v-model="formData.amount"
                  type="number"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div class="sm:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                  v-model="formData.description"
                  rows="3"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              ></textarea>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Picture URLs</label>
            <div class="space-y-2">
              <div v-for="(url, index) in formData.pictureUrls" :key="index" class="flex gap-2">
                <input
                    v-model="formData.pictureUrls[index]"
                    type="text"
                    placeholder="Image URL"
                    class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                >
                <button
                    v-if="formData.pictureUrls.length > 1"
                    @click="removePictureUrl(index)"
                    type="button"
                    class="px-3 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-lg"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              <button
                  @click="addPictureUrl"
                  type="button"
                  class="w-full px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-600 hover:bg-gray-50"
              >
                + Add another image
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Colors</label>
            <div class="space-y-2">
              <div v-for="(color, index) in formData.colors" :key="index" class="flex gap-2">
                <input
                    v-model="color.name"
                    type="text"
                    placeholder="Color name"
                    class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                >
                <button
                    v-if="formData.colors.length > 1"
                    @click="removeColor(index)"
                    type="button"
                    class="px-3 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-lg"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              <button
                  @click="addColor"
                  type="button"
                  class="w-full px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-600 hover:bg-gray-50"
              >
                + Add another color
              </button>
            </div>
          </div>

          <div class="sticky bottom-0 bg-white pt-4 border-t border-gray-200 flex flex-col sm:flex-row gap-2 sm:gap-3">
            <button
                @click="closeModals"
                type="button"
                class="flex-1 px-4 py-3 sm:py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
            >
              Cancel
            </button>
            <button
                type="submit"
                class="flex-1 px-4 py-3 sm:py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
            >
              Add Product
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Mobile-Optimized Edit Product Modal -->
    <div v-if="isEditModalOpen"
         class="fixed inset-0 bg-black bg-opacity-50 flex items-end sm:items-center justify-center z-50"
         @click.self="closeModals">
      <div
          class="bg-white rounded-t-lg sm:rounded-lg shadow-xl w-full sm:max-w-2xl sm:mx-4 max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white px-4 sm:px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">Edit Product</h3>
          <button
              @click="closeModals"
              class="p-2 hover:bg-gray-100 rounded-lg"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        <form @submit.prevent="updateProduct" class="p-4 sm:p-6 space-y-4 sm:space-y-6">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
            <div class="sm:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Name</label>
              <input
                  v-model="formData.name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Model</label>
              <input
                  v-model="formData.model"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Event</label>
              <select
                  v-model="formData.event"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="none">None</option>
                <option value="sale">Sale</option>
                <option value="featured">Featured</option>
                <option value="new">New</option>
                <option value="bestseller">Bestseller</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Price ($)</label>
              <input
                  v-model="formData.price"
                  type="number"
                  step="0.01"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Stock</label>
              <input
                  v-model="formData.amount"
                  type="number"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
            </div>
            <div class="sm:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                  v-model="formData.description"
                  rows="3"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              ></textarea>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Picture URLs</label>
            <div class="space-y-2">
              <div v-for="(url, index) in formData.pictureUrls" :key="index" class="flex gap-2">
                <input
                    v-model="formData.pictureUrls[index]"
                    type="text"
                    placeholder="Image URL"
                    class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                >
                <button
                    v-if="formData.pictureUrls.length > 1"
                    @click="removePictureUrl(index)"
                    type="button"
                    class="px-3 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-lg"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              <button
                  @click="addPictureUrl"
                  type="button"
                  class="w-full px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-600 hover:bg-gray-50"
              >
                + Add another image
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Colors</label>
            <div class="space-y-2">
              <div v-for="(color, index) in formData.colors" :key="index" class="flex gap-2">
                <input
                    v-model="color.name"
                    type="text"
                    placeholder="Color name"
                    class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                >
                <button
                    v-if="formData.colors.length > 1"
                    @click="removeColor(index)"
                    type="button"
                    class="px-3 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-lg"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              <button
                  @click="addColor"
                  type="button"
                  class="w-full px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-600 hover:bg-gray-50"
              >
                + Add another color
              </button>
            </div>
          </div>

          <div class="sticky bottom-0 bg-white pt-4 border-t border-gray-200 flex flex-col sm:flex-row gap-2 sm:gap-3">
            <button
                @click="closeModals"
                type="button"
                class="flex-1 px-4 py-3 sm:py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
            >
              Cancel
            </button>
            <button
                type="submit"
                class="flex-1 px-4 py-3 sm:py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
            >
              Update Product
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Mobile-Optimized Delete Confirmation Modal -->
    <div v-if="isDeleteModalOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
         @click.self="closeModals">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="px-4 sm:px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Delete Product</h3>
        </div>
        <div class="p-4 sm:p-6">
          <p class="text-gray-600 mb-4">
            Are you sure you want to delete "{{ selectedProduct?.name }}"? This action cannot be undone.
          </p>
          <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
            <button
                @click="closeModals"
                class="flex-1 px-4 py-3 sm:py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
            >
              Cancel
            </button>
            <button
                @click="deleteProduct"
                class="flex-1 px-4 py-3 sm:py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 font-medium"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar for modals */
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* Line clamp utility for text truncation */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Prevent scroll when modals are open */
body.modal-open {
  overflow: hidden;
}

/* Improve touch targets on mobile */
@media (max-width: 768px) {
  button, input, select, textarea {
    min-height: 44px;
  }

  .touch-friendly {
    min-height: 44px;
    min-width: 44px;
  }
}

/* Ensure proper spacing on very small screens */
@media (max-width: 375px) {
  .px-4 {
    padding-left: 1rem;
    padding-right: 1rem;
  }
}
</style>