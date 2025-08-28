<script setup>
import Wrapper from "@/components/Wrapper.vue";
import {ArrowRight, Menu, Search, ShoppingBag, User} from 'lucide-vue-next';
import SheetGeneral from "@/components/SheetGeneral.vue";
import {useRoute} from "vue-router";
import {computed, ref, watch} from "vue";
import {cn} from "@/lib/utils.js";
import {Button} from "@/components/ui/button/index.js";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.ts";
import {useParams} from "@/composables/useParams.js";
import {useAuthStore} from "@/stores/useAuth.js";
import Profile from "@/components/Profile.vue";
import AuthComponent from "@/components/auth/AuthComponent.vue";
import RegistrationButton from "@/components/auth/RegistrationButton.vue";
import LoginButton from "@/components/auth/LoginButton.vue";

const smallScreenSheet = useScreenSheetStore()
const auth = useAuthStore()
const route = useRoute()
const {updateProductName} = useParams()

const isOnAboutPage = computed(() => route.path === "/about");
const isOnContactPage = computed(() => route.path === "/contact");
const isOnProductsPage = computed(() => route.path === "/products");
const isOnProfilePage = computed(() => route.path.startsWith("/profile"));

const productName = ref("")
const timeout = ref(null)
const updatedRegisterState = computed(() => smallScreenSheet.isDialogOpen("RegisterDialog"))
const updatedLoginState = computed(() => smallScreenSheet.isDialogOpen("LoginDialog"))

const handleChange = (value) => {
  productName.value = value;
}

const handleSearch = () => {
  smallScreenSheet.setOpenSheet(false, {name: 'SmallScreenSearch'})
  updateProductName(productName.value);
}

watch(updatedLoginState, () => {
  auth.error = ""
})
watch(updatedRegisterState, () => {
  auth.error = ""
})


watch(productName, (newValue) => {
  if (window.innerWidth >= 1024) {
    if (timeout.value) clearTimeout(timeout.value)

    timeout.value = setTimeout(() => {
      updateProductName(newValue);
    }, 500)
  }
})

</script>

<template>
  <div class="z-50 w-full backdrop-blur-sm bg-white/70 fixed top-0">
    <Wrapper>
      <div
          :class="cn(isOnProductsPage ? 'grid lg:grid-cols-[minmax(200px,328px)_1fr_minmax(100px,auto)] sm:grid-cols-[minmax(80px,132px)_1fr_minmax(80px,auto)]' : 'grid grid-cols-[1fr_3fr_1fr]')">
        <SheetGeneral :is-open="smallScreenSheet.isSheetOpen('SmallScreenMenu')" trigger-class="sm:hidden" side="left"
                      title="FUMI">
          <template #trigger>
            <Menu @click="smallScreenSheet.setOpenSheet(true, {name: 'SmallScreenMenu'})"
                  class="order-1 block sm:hidden cursor-pointer col-start-1"/>
          </template>
          <template #content>
            <ul class="p-4">
              <li>
                <div class=" text-gray-700 hover:text-black">
                  <router-link @click="smallScreenSheet.setOpenSheet(false, {name: 'SmallScreenMenu'})" to="/products"
                               class="text-2xl flex justify-between w-full items-center">
                    <p class="text-2xl ">Products</p>
                    <ArrowRight/>
                  </router-link>
                </div>
              </li>
              <li>
                <div class="text-gray-700 hover:text-black">
                  <router-link @click="smallScreenSheet.setOpenSheet(false, {name: 'SmallScreenMenu'})" to="/about"
                               class="text-2xl flex justify-between w-full items-center">
                    About Us
                    <ArrowRight/>
                  </router-link>
                </div>
              </li>
              <li>
                <div class="text-gray-700 hover:text-black">
                  <router-link @click="smallScreenSheet.setOpenSheet(false, {name: 'SmallScreenMenu'})" to="/contact"
                               class="text-2xl flex w-full justify-between items-center">
                    Contact Us
                    <ArrowRight/>
                  </router-link>
                </div>
              </li>
            </ul>
            <ul class="p-4">
              <li v-if="auth.isAuthenticated">
                <router-link @click="smallScreenSheet.setOpenSheet(false, {name: 'SmallScreenMenu'})"
                             class="flex items-center gap-2 text-2xl text-gray-700 hover:text-black"
                             :to="`/profile/${auth.user.id}`">
                  <User/>
                  Profile
                </router-link>
              </li>
              <li v-else class="flex items-center gap-2 text-2xl text-gray-700 hover:text-black">
                <RegistrationButton class-name="flex-1 cursor-pointer"/>
                <LoginButton class-name="flex-1 cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]"/>
              </li>
            </ul>
          </template>
        </SheetGeneral>
        <router-link to="/"
                     class="order-2 sm:order-1 col-start-2 sm:col-start-1 flex items-center justify-center sm:justify-start">
          <h1 class="font-semibold text-base sm:text-lg">FUMI</h1>
        </router-link>
        <ul class="sm:order-2 col-start-2 hidden sm:flex sm:items-center sm:justify-center gap-3">
          <li :class="cn('text-gray-700 hover:text-black', isOnProductsPage && 'text-black')">
            <router-link to="/products">
              Products
            </router-link>
          </li>
          <li :class="cn('text-gray-700 hover:text-black', isOnAboutPage && 'text-black')">
            <router-link to="/about">
              About Us
            </router-link>
          </li>
          <li :class="cn('text-gray-700 hover:text-black', isOnContactPage && 'text-black')">
            <router-link to="/contact">
              Contact Us
            </router-link>
          </li>
        </ul>
        <ul class="order-3 col-start-3 flex items-center justify-end gap-3">
          <li v-if="isOnProductsPage" class="relative items-center hidden lg:flex">
            <input v-model="productName" name="search" class="pl-7 border rounded-sm px-2 py-1 bg-white" type="text"
                   placeholder="Search..."/>
            <Search class="text-gray-600 absolute top-[5px] left-1"/>
          </li>
          <SheetGeneral :is-open="smallScreenSheet.isSheetOpen('SmallScreenSearch')" side="top">
            <template #trigger>
              <li v-if="isOnProductsPage" @click="smallScreenSheet.setOpenSheet(true, {name: 'SmallScreenSearch'})"
                  class="lg:hidden block">
                <Search></Search>
              </li>
            </template>
            <template #content>
              <div class="relative items-center w-full flex p-5 pt-0 gap-3">
                <input @change="(e) => handleChange(e.target.value)" name="search"
                       class="pl-7 border w-full rounded-sm px-2 py-1" type="text"
                       placeholder="Search..."/>
                <Search class="text-gray-600 absolute top-[5px] left-6"/>
                <Button @click="handleSearch" class="bg-[#c9a275] h-[34px] cursor-pointer hover:bg-[#dbb384]">Go
                </Button>
              </div>
            </template>
          </SheetGeneral>
          <SheetGeneral :is-open="smallScreenSheet.isSheetOpen('ScreenProfile')" title="">
            <template #trigger>
              <li v-if="!isOnProfilePage" class="hidden sm:block"
                  @click="smallScreenSheet.setOpenSheet(true, {name: 'ScreenProfile'})">
                <User></User>
              </li>
            </template>
            <template #content>
              <div class="flex items-center justify-center h-full">
                <Profile v-if="auth.isAuthenticated"/>
                <AuthComponent v-else/>
              </div>
            </template>
          </SheetGeneral>
          <SheetGeneral :is-open="smallScreenSheet.isSheetOpen('ScreenCart')" title="Your Cart">
            <template #trigger>
              <li @click="smallScreenSheet.setOpenSheet(true, {name: 'ScreenCart'})">
                <ShoppingBag></ShoppingBag>
              </li>
            </template>
            <template #content>
              <h1>Content!</h1>
            </template>
          </SheetGeneral>
        </ul>
      </div>
    </Wrapper>
  </div>
</template>