<template>
  <div style="font-family: sans-serif; max-width: 800px; margin: 20px auto;">
    <h1>Mini CRM Automation UI</h1>

    <div style="border: 1px solid #ddd; padding: 20px; border-radius: 5px;">
      <h3>Submit Deposit Event</h3>
      <input v-model="player_id" placeholder="Player ID" style="margin-right: 10px; padding: 5px;"/>
      <input v-model.number="amount" type="number" placeholder="Amount" style="margin-right: 10px; padding: 5px;"/>
      <button @click="submit" style="padding: 5px 15px; cursor: pointer;">Submit</button>
    </div>

    <h3 style="margin-top: 40px;">Campaign Logs (from MongoDB)</h3>
    <table border="1" style="width: 100%; border-collapse: collapse;">
      <thead>
        <tr style="background: #f4f4f4;">
          <th>Player ID</th><th>Action</th><th>Amount</th><th>Timestamp</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="log in logs" :key="log.event_id">
          <td>{{ log.player_id }}</td>
          <td style="color: green; font-weight: bold;">{{ log.action }}</td>
          <td>{{ log.amount }}</td>
          <td>{{ new Date(log.timestamp).toLocaleString() }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  data: () => ({ player_id: '', amount: 0, logs: [] }),
  methods: {
    async submit() {
      await fetch('http://localhost:8080/ingest', {
        method: 'POST',
        body: JSON.stringify({ player_id: this.player_id, amount: this.amount })
      });
      this.fetchLogs();
    },
    async fetchLogs() {
      const res = await fetch('http://localhost:8080/logs');
      this.logs = await res.json();
    }
  },
  mounted() { this.fetchLogs(); }
}
</script>