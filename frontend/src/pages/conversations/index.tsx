import { useEffect, useState } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { conversations as conversationsAPI } from '@/lib/api';
import DashboardLayout from '@/components/DashboardLayout';

interface Conversation {
  id: string;
  patient_phone: string;
  status: string;
  intent: string;
  summary: string;
  is_resolved: boolean;
  created_at: string;
}

export default function Conversations() {
  const router = useRouter();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadConversations();
  }, []);

  const loadConversations = async () => {
    try {
      const response = await conversationsAPI.list();
      setConversations(response.data);
    } catch (error) {
      console.error('Failed to load conversations:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Head>
        <title>Conversations - RecallFlow</title>
      </Head>

      <DashboardLayout>
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Conversations</h1>
            <p className="text-gray-600 mt-1">Manage all patient conversations</p>
          </div>

          {loading ? (
            <div className="text-center py-12">Loading...</div>
          ) : conversations.length === 0 ? (
            <div className="card text-center py-12">
              <p className="text-gray-500">No conversations yet</p>
            </div>
          ) : (
            <div className="space-y-4">
              {conversations.map((conv) => (
                <div
                  key={conv.id}
                  className="card hover:shadow-lg transition-shadow cursor-pointer"
                  onClick={() => router.push(`/conversations/${conv.id}`)}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-4">
                        <div>
                          <div className="font-medium text-gray-900">{conv.patient_phone}</div>
                          <div className="text-sm text-gray-500">
                            {new Date(conv.created_at).toLocaleString()}
                          </div>
                        </div>
                        <span className="inline-flex px-3 py-1 text-sm font-medium rounded bg-blue-100 text-blue-800">
                          {conv.intent || 'unknown'}
                        </span>
                        <span className={`inline-flex px-3 py-1 text-sm font-medium rounded ${
                          conv.status === 'active' ? 'bg-green-100 text-green-800' :
                          conv.status === 'resolved' ? 'bg-gray-100 text-gray-800' :
                          'bg-yellow-100 text-yellow-800'
                        }`}>
                          {conv.status}
                        </span>
                      </div>
                      {conv.summary && (
                        <p className="mt-2 text-sm text-gray-600">{conv.summary}</p>
                      )}
                    </div>
                    <div className="text-primary-600 font-medium">→</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </DashboardLayout>
    </>
  );
}
