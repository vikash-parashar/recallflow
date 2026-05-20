import { useEffect, useState } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { dashboard, conversations as conversationsAPI } from '@/lib/api';
import DashboardLayout from '@/components/DashboardLayout';

interface Stats {
  total_calls: number;
  missed_calls: number;
  recovered_leads: number;
  active_conversations: number;
  response_rate: number;
  estimated_revenue: number;
}

interface Conversation {
  id: string;
  patient_phone: string;
  status: string;
  intent: string;
  created_at: string;
}

export default function Dashboard() {
  const router = useRouter();
  const [stats, setStats] = useState<Stats | null>(null);
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      router.push('/login');
      return;
    }

    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [statsRes, conversationsRes] = await Promise.all([
        dashboard.getStats(),
        conversationsAPI.list(),
      ]);
      
      setStats(statsRes.data);
      setConversations(conversationsRes.data.slice(0, 5)); // Latest 5
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-gray-500">Loading...</div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <>
      <Head>
        <title>Dashboard - RecallFlow</title>
      </Head>

      <DashboardLayout>
        <div className="space-y-6">
          {/* Header */}
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <p className="text-gray-600 mt-1">Welcome back! Here's your clinic overview.</p>
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div className="card">
              <div className="text-sm font-medium text-gray-600">Total Calls (30d)</div>
              <div className="text-3xl font-bold text-gray-900 mt-2">{stats?.total_calls || 0}</div>
            </div>

            <div className="card">
              <div className="text-sm font-medium text-gray-600">Missed Calls</div>
              <div className="text-3xl font-bold text-red-600 mt-2">{stats?.missed_calls || 0}</div>
            </div>

            <div className="card">
              <div className="text-sm font-medium text-gray-600">Recovered Leads</div>
              <div className="text-3xl font-bold text-green-600 mt-2">{stats?.recovered_leads || 0}</div>
            </div>

            <div className="card">
              <div className="text-sm font-medium text-gray-600">Active Conversations</div>
              <div className="text-3xl font-bold text-blue-600 mt-2">{stats?.active_conversations || 0}</div>
            </div>

            <div className="card">
              <div className="text-sm font-medium text-gray-600">Response Rate</div>
              <div className="text-3xl font-bold text-gray-900 mt-2">{stats?.response_rate.toFixed(1) || 0}%</div>
            </div>

            <div className="card">
              <div className="text-sm font-medium text-gray-600">Estimated Revenue</div>
              <div className="text-3xl font-bold text-green-600 mt-2">${stats?.estimated_revenue.toLocaleString() || 0}</div>
            </div>
          </div>

          {/* Recent Conversations */}
          <div className="card">
            <h2 className="text-xl font-semibold mb-4">Recent Conversations</h2>
            {conversations.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No conversations yet</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Patient</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Intent</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Time</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {conversations.map((conv) => (
                      <tr key={conv.id} className="hover:bg-gray-50">
                        <td className="px-4 py-3 text-sm">{conv.patient_phone}</td>
                        <td className="px-4 py-3 text-sm">
                          <span className="inline-flex px-2 py-1 text-xs font-medium rounded bg-blue-100 text-blue-800">
                            {conv.intent || 'unknown'}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-sm">
                          <span className={`inline-flex px-2 py-1 text-xs font-medium rounded ${
                            conv.status === 'active' ? 'bg-green-100 text-green-800' :
                            conv.status === 'resolved' ? 'bg-gray-100 text-gray-800' :
                            'bg-yellow-100 text-yellow-800'
                          }`}>
                            {conv.status}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-600">
                          {new Date(conv.created_at).toLocaleString()}
                        </td>
                        <td className="px-4 py-3 text-sm">
                          <button
                            onClick={() => router.push(`/conversations/${conv.id}`)}
                            className="text-primary-600 hover:text-primary-800"
                          >
                            View
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </DashboardLayout>
    </>
  );
}
