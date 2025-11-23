import Link from 'next/link';

export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <h1 className="text-4xl font-bold text-center sm:text-left">
          Fish Auction System
        </h1>
        <p className="text-lg text-gray-600">
          Welcome to the digital auction floor.
        </p>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 w-full">
          <Link href="/admin" className="block p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
            <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Admin</h5>
            <p className="font-normal text-gray-700 dark:text-gray-400">Register Fishermen, Buyers, and Items.</p>
          </Link>

          <Link href="/auction" className="block p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
            <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Auction</h5>
            <p className="font-normal text-gray-700 dark:text-gray-400">View items and place bids.</p>
          </Link>

          <Link href="/invoice" className="block p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
            <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Invoices</h5>
            <p className="font-normal text-gray-700 dark:text-gray-400">View billing information.</p>
          </Link>
        </div>
      </main>
    </div>
  );
}
