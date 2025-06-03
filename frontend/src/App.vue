<template>
  <div>
    <nav>
      <a href="/">Home</a> | <a href="/login">Login</a>
    </nav>
    <router-view />
  </div>
</template>

<style>
nav {
  padding: 1em;
  background: #f0f0f0;
}
a {
  margin: 0 0.5em;
  text-decoration: none;
  color: #333;
}
</style>

