<script setup>
import RegistrationButton from "@/components/auth/RegistrationButton.vue";
import LoginButton from "@/components/auth/LoginButton.vue";
</script>

<template>
  <div class="flex flex-col justify-center items-center gap-12">
    <h1 class="font-semibold text-4xl">FUMI</h1>
    <div class="flex justify-around items-center h-full gap-3">
      <RegistrationButton class-name="text-lg cursor-pointer"/>
      <LoginButton class-name="text-lg cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]"/>
    </div>
  </div>
</template>