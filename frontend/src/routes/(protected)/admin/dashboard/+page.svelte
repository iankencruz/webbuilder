<script lang="ts">
  import { stats, transactions } from '$lib/data/dashboard';
</script>

<div class="space-y-6">
  <!-- Page title -->
  <div>
    <h1 class="text-xl font-semibold text-foreground">Dashboard</h1>
  </div>

  <!-- Stats -->
  <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
    {#each stats as stat}
      <div class="bg-card rounded-lg border border-border p-5 shadow-sm">
        <p class="text-sm text-muted-foreground">{stat.label}</p>
        <p class="text-2xl font-semibold text-card-foreground mt-1">
          {stat.value}
        </p>
        <p
          class="text-xs mt-1 {stat.up ? 'text-green-600' : 'text-destructive'}"
        >
          {stat.change} from last month
        </p>
      </div>
    {/each}
  </div>

  <!-- Transactions -->
  <div class="bg-card rounded-lg border border-border shadow-sm">
    <div
      class="flex items-center justify-between px-5 py-4 border-b border-border"
    >
      <h2 class="text-sm font-semibold text-card-foreground">
        Recent Transactions
      </h2>
      <button class="text-xs text-primary hover:underline">View all</button>
    </div>
    <table class="w-full text-sm">
      <thead>
        <tr
          class="text-left text-xs text-muted-foreground border-b border-border"
        >
          <th class="px-5 py-3 font-medium">Name</th>
          <th class="px-5 py-3 font-medium">Description</th>
          <th class="px-5 py-3 font-medium">Amount</th>
          <th class="px-5 py-3 font-medium hidden sm:table-cell">Date</th>
          <th class="px-5 py-3 font-medium">Status</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-border">
        {#each transactions as tx}
          <tr class="hover:bg-accent transition-colors">
            <td class="px-5 py-3.5">
              <div class="flex items-center gap-2.5">
                <div
                  class="w-7 h-7 rounded-full bg-muted flex items-center justify-center shrink-0"
                >
                  <span class="text-muted-foreground text-xs font-medium"
                    >{tx.name[0]}</span
                  >
                </div>
                <span class="font-medium text-foreground">{tx.name}</span>
              </div>
            </td>
            <td class="px-5 py-3.5 text-muted-foreground">{tx.desc}</td>
            <td
              class="px-5 py-3.5 font-medium {tx.amount.startsWith('+')
                ? 'text-green-600'
                : 'text-destructive'}"
            >
              {tx.amount}
            </td>
            <td class="px-5 py-3.5 text-muted-foreground hidden sm:table-cell"
              >{tx.date}</td
            >
            <td class="px-5 py-3.5">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium
                {tx.status === 'Credit'
                  ? 'bg-green-50 text-green-700'
                  : 'bg-destructive/10 text-destructive'}"
              >
                {tx.status}
              </span>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
