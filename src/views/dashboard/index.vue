<template>
  <div class="dashboard-container">
    <div class="dashboard-text">name:{{ name }}</div>
    <div class="dashboard-text">roles:<span v-for="role in roles" :key="role">{{ role }}</span></div>
    <table>
      <th>id</th>
      <th>name</th>
      <tr v-for="(t,i) in ta" :key="i">
        <td>{{ t.id }}</td>
        <td>{{ t.name }}</td>
      </tr>
      <button @click="fetch">test</button>
    </table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { test } from '@/api/test'

export default {
  data() {
    return {
      ta: [{ 'id': 123, 'name': 222 }, { 'id': 1233, 'name': 344 }]
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'roles'
    ])
  },
  methods: {
    fetch: function() {
      test('token')
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
