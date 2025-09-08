"use client"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { getFIP } from "@/data/loaders/fips"
import { getUpgrade } from "@/data/loaders/upgrades"
import upgradesData from "@/data/upgrades.json"
import { motion } from "framer-motion"
import { ArrowLeft, Code, DollarSign, Server, Star, TrendingUp, User } from "lucide-react"
import { useParams, useRouter } from "next/navigation"
import { useEffect, useState } from "react"
import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts"

const { metadata } = upgradesData
const { statusColors, fallbacks } = metadata

const categoryColors = {
  Storage: "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300",
  Retrieval: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
  Security: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  Performance: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
  Economic: "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300",
}

const isFip0077 = (fip: any) => {
  const id = (fip?.id || "").toString().toUpperCase().replace(/\s+/g, "")
  return id === "FIP-0077" || id === "FIP0077" || fip?.number === 77
}

interface MinerCountData {
  date: string;
  count: number;
}

interface FIPImpactModalProps {
  fip: any;
  isOpen: boolean;
  onClose: () => void;
}

function FIPImpactModal({ fip, isOpen, onClose }: FIPImpactModalProps) {
  const [minerData, setMinerData] = useState<MinerCountData[]>([]);
  const [timeRange, setTimeRange] = useState<'7d' | '30d'>('7d');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch data when timeRange changes
  useEffect(() => {
    if (!isOpen || !fip) return;

    const fetchMinerData = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await fetch(`/api/miner-count?range=${timeRange}`);
        if (!response.ok) {
          throw new Error('Failed to fetch miner data');
        }
        const data = await response.json();
        console.log(data);
        setMinerData(data);
      } catch (err) {
        setError('Error fetching data. Please try again.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchMinerData();
  }, [timeRange, isOpen, fip]);

  if (!fip) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Badge className={categoryColors.Storage}>{fip.category || 'General'}</Badge>
            {fip.id}: {fip.title}
          </DialogTitle>
          <p className="text-muted-foreground">{fip.description}</p>
        </DialogHeader>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <TrendingUp className="h-5 w-5" />
                Filecoin Miner Count Trends
              </CardTitle>
              <p className="text-sm text-muted-foreground">
                Historical and projected changes in Filecoin network miner participation
              </p>
              <div className="mt-2">
                <select
                  value={timeRange}
                  onChange={(e) => setTimeRange(e.target.value as '7d' | '30d')}
                  className="w-32 rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-primary"
                >
                  <option value="7d">7 Days</option>
                  <option value="30d">30 Days</option>
                </select>
              </div>
              {error && (
                <p className="text-xs text-red-500 font-medium mt-1">{error}</p>
              )}
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="h-80 w-full flex items-center justify-center">
                  <p>Loading...</p>
                </div>
              ) : (
                <div className="h-80 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={minerData}>
                      <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                      <XAxis
                        dataKey="date"
                        tick={{ fontSize: 12 }}
                        tickLine={{ stroke: "#666" }}
                      />
                      <YAxis
                        tick={{ fontSize: 12 }}
                        tickLine={{ stroke: "#666" }}
                        domain={[
                          (dataMin: number) => Math.max(0, dataMin - 5),
                          (dataMax: number) => dataMax + 5,
                        ]}
                      />
                      <Tooltip
                        contentStyle={{
                          backgroundColor: 'hsl(var(--background))',
                          border: '1px solid hsl(var(--border))',
                          borderRadius: '6px',
                          fontSize: '12px',
                        }}
                        labelFormatter={(label) => `Date: ${label}`}
                        formatter={(value, name) => [
                          value.toLocaleString(),
                          name === 'count' ? 'Total Miners' : 'Active Miners',
                        ]}
                      />
                      <Line
                        type="monotone"
                        dataKey="count"
                        stroke="#8884d8"
                        strokeWidth={2}
                        dot={{ fill: '#8884d8', stroke: '#8884d8', strokeWidth: 2, r: 4 }}
                        activeDot={{ r: 6, stroke: '#8884d8', strokeWidth: 2 }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </DialogContent>
    </Dialog>
  );
}


const loadUpgrade = async (upgradeId: string) => {
  try {
    const upgrade = getUpgrade(upgradeId)
    if (!upgrade) {
      throw new Error(`Upgrade ${upgradeId} not found`)
    }
    return upgrade
  } catch (error) {
    console.error(`Failed to load upgrade ${upgradeId}:`, error)
    return null
  }
}

const loadFIP = async (fipId: string) => {
  try {
    const fip = getFIP(fipId)
    if (!fip) {
      throw new Error(`FIP ${fipId} not found`)
    }
    return fip
  } catch (error) {
    console.error(`Failed to load FIP ${fipId}:`, error)
    return null
  }
}

export default function UpgradeDetails() {
  const router = useRouter()
  const params = useParams()
  const upgradeId = params.id as string
  const [selectedFIP, setSelectedFIP] = useState<any | null>(null)
  const [activeFIPId, setActiveFIPId] = useState<string | null>(null)
  const [upgrade, setUpgrade] = useState<any>(null)
  const [fips, setFips] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const loadUpgradeData = async () => {
      setLoading(true)

      const upgradeData = await loadUpgrade(upgradeId)
      if (!upgradeData) {
        setLoading(false)
        return
      }

      setUpgrade(upgradeData)

      const loadedFips = []
      for (const fipId of upgradeData.fipIds || []) {
        const fip = await loadFIP(fipId)
        if (fip) {
          loadedFips.push(fip)
        }
      }

      setFips(loadedFips)
      setLoading(false)
    }

    loadUpgradeData()
  }, [upgradeId])

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading upgrade details...</p>
        </div>
      </div>
    )
  }

  if (!upgrade) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Upgrade Not Found</h1>
          <Button onClick={() => router.push("/")}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Upgrades
          </Button>
        </div>
      </div>
    )
  }

  const safeUpgrade = {
    ...upgrade,
    status: upgrade.status || fallbacks.defaultStatus,
    notes: upgrade.notes || fallbacks.defaultNotes,
    releaseTag: upgrade.releaseTag || fallbacks.defaultReleaseTag,
    specs: upgrade.specs || fallbacks.emptySpecs,
    fips: fips,
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-3">
            <Button variant="link" className="p-0 h-auto text-sm" onClick={() => router.push("/")}>
              All Network Upgrades
            </Button>
            <span>/</span>
            <span>{safeUpgrade.name}</span>
          </div>
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-balance mb-2">{safeUpgrade.name}</h1>
              <p className="text-muted-foreground">{safeUpgrade.notes}</p>
            </div>
            <Badge className={statusColors[safeUpgrade.status as keyof typeof statusColors] || statusColors.Upcoming}>
              {safeUpgrade.status}
            </Badge>
          </div>
        </div>
      </header>

      <div className="border-b bg-muted/20">
        <div className="container mx-auto px-4 py-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div>
              <div className="text-2xl font-bold text-muted-foreground">0</div>
              <div className="text-sm text-muted-foreground">Activated</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-muted-foreground">0</div>
              <div className="text-sm text-muted-foreground">Scheduled</div>
            </div>
            <div>
              <div className="text-2xl font-bold">{safeUpgrade.fips.length}</div>
              <div className="text-sm text-muted-foreground">FIPs</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-muted-foreground">0</div>
              <div className="text-sm text-muted-foreground">Declined or Postponed</div>
            </div>
          </div>
        </div>
      </div>

      <main className="container mx-auto px-4 py-8">
        <div className="flex gap-8">
          <div className="w-64 flex-shrink-0">
            <div className="sticky top-8">
              <h3 className="font-semibold text-sm uppercase tracking-wide text-muted-foreground mb-4">Contents</h3>
              <nav className="space-y-1">
                <a href="#overview" className="block text-sm hover:text-primary transition-colors py-1">
                  Overview
                </a>
                <a href="#fips" className="block text-sm text-primary font-medium py-1">
                  FIPs
                </a>
                {safeUpgrade.fips.map((fip: any) => (
                  <a
                    key={fip.id}
                    href={`#${fip.id}`}
                    className={`group block text-sm transition-colors pl-4 py-1 font-mono ${activeFIPId === fip.id
                      ? "text-primary font-medium bg-primary/10 rounded-md px-2"
                      : "text-muted-foreground hover:text-primary"
                      }`}
                    onClick={() => setActiveFIPId(fip.id)}
                    title={fip.title}
                  >
                    <span className="flex items-center gap-2 min-w-0">
                      {isFip0077(fip) && (
                        <Star className="h-3.5 w-3.5 text-amber-500 flex-shrink-0" aria-label="Important FIP" />
                      )}
                      <span className="flex-shrink-0">{fip.id}: </span>
                      <span
                        className={`truncate flex-1 min-w-0 ${activeFIPId === fip.id
                          ? "text-primary"
                          : "text-muted-foreground group-hover:text-primary"
                          }`}
                      >
                        {fip.title}
                      </span>
                    </span>
                  </a>
                ))}
              </nav>
            </div>
          </div>

          <div className="flex-1 max-w-4xl">
            <section id="overview" className="mb-12">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
                <div className="bg-muted/30 rounded-lg p-4">
                  <div className="text-sm text-muted-foreground mb-1">Network Version</div>
                  <div className="text-lg font-mono">{safeUpgrade.networkVersion}</div>
                </div>
                <div className="bg-muted/30 rounded-lg p-4">
                  <div className="text-sm text-muted-foreground mb-1">Target Date</div>
                  <div className="text-lg">
                    {safeUpgrade.timeTarget
                      ? new Date(safeUpgrade.timeTarget).toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      })
                      : "TBD"}
                  </div>
                </div>
                <div className="bg-muted/30 rounded-lg p-4">
                  <div className="text-sm text-muted-foreground mb-1">Release Tag</div>
                  <div className="text-lg font-mono">{safeUpgrade.releaseTag}</div>
                </div>
              </div>
            </section>

            <section id="fips" className="mb-12">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-2xl font-bold">FIPs</h2>
                <span className="text-sm text-muted-foreground">({safeUpgrade.fips.length})</span>
              </div>

              <p className="text-muted-foreground mb-8">The following FIPs are part of this upgrade.</p>

              <div className="space-y-6">
                {safeUpgrade.fips.length > 0 ? (
                  safeUpgrade.fips.map((fip: any, index: number) => (
                    <motion.div
                      key={fip.id}
                      id={fip.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.3, delay: index * 0.1 }}
                      className={`border rounded-lg p-6 hover:shadow-sm transition-all ${activeFIPId === fip.id ? "border-primary shadow-sm bg-primary/5" : ""
                        }`}
                    >
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-3">
                            <Badge variant="outline" className="font-mono">
                              {fip.id}
                            </Badge>
                            {isFip0077(fip) && <Star className="h-4 w-4 text-amber-500" aria-label="Important FIP" />}
                            {fip.category && (
                              <Badge
                                className={
                                  categoryColors[fip.category as keyof typeof categoryColors] || categoryColors.Storage
                                }
                              >
                                {fip.category}
                              </Badge>
                            )}
                          </div>
                          <h3 className="text-xl font-semibold mb-2">{fip.title}</h3>
                          {fip.description && <p className="text-sm text-muted-foreground mb-4">{fip.description}</p>}
                        </div>
                      </div>

                      {fip.impacts && (
                        <div className="border-t pt-6">
                          <div className="space-y-6">
                            {fip.impacts.storageProviders && (
                              <div>
                                <h5 className="font-medium text-blue-700 dark:text-blue-400 mb-3 flex items-center gap-2">
                                  <Server className="h-4 w-4" />
                                  Storage Providers
                                </h5>
                                <ul className="ml-6 space-y-1">
                                  {fip.impacts.storageProviders.length > 0 ? (
                                    fip.impacts.storageProviders.map((impact: string, idx: number) => (
                                      <li key={idx} className="text-sm text-muted-foreground flex items-start gap-2">
                                        <span className="mt-1 text-xs">•</span>
                                        {impact}
                                      </li>
                                    ))
                                  ) : (
                                    <li className="text-sm text-muted-foreground">No impact specified yet.</li>
                                  )}
                                </ul>
                              </div>
                            )}

                            {fip.impacts.clients && (
                              <div>
                                <h5 className="font-medium text-orange-700 dark:text-orange-400 mb-3 flex items-center gap-2">
                                  <User className="h-4 w-4" />
                                  Clients
                                </h5>
                                <ul className="ml-6 space-y-1">
                                  {fip.impacts.clients.length > 0 ? (
                                    fip.impacts.clients.map((impact: string, idx: number) => (
                                      <li key={idx} className="text-sm text-muted-foreground flex items-start gap-2">
                                        <span className="mt-1 text-xs">•</span>
                                        {impact}
                                      </li>
                                    ))
                                  ) : (
                                    <li className="text-sm text-muted-foreground">No impact specified yet.</li>
                                  )}
                                </ul>
                              </div>
                            )}

                            {fip.impacts.tokenHolders && (
                              <div>
                                <h5 className="font-medium text-purple-700 dark:text-purple-400 mb-3 flex items-center gap-2">
                                  <DollarSign className="h-4 w-4" />
                                  Token Holders
                                </h5>
                                <ul className="ml-6 space-y-1">
                                  {fip.impacts.tokenHolders.length > 0 ? (
                                    fip.impacts.tokenHolders.map((impact: string, idx: number) => (
                                      <li key={idx} className="text-sm text-muted-foreground flex items-start gap-2">
                                        <span className="mt-1 text-xs">•</span>
                                        {impact}
                                      </li>
                                    ))
                                  ) : (
                                    <li className="text-sm text-muted-foreground">No impact specified yet.</li>
                                  )}
                                </ul>
                              </div>
                            )}

                            {fip.impacts.applicationDevelopers && (
                              <div>
                                <h5 className="font-medium text-green-700 dark:text-green-400 mb-3 flex items-center gap-2">
                                  <Code className="h-4 w-4" />
                                  Application Developers
                                </h5>
                                <ul className="ml-6 space-y-1">
                                  {fip.impacts.applicationDevelopers.length > 0 ? (
                                    fip.impacts.applicationDevelopers.map((impact: string, idx: number) => (
                                      <li key={idx} className="text-sm text-muted-foreground flex items-start gap-2">
                                        <span className="mt-1 text-xs">•</span>
                                        {impact}
                                      </li>
                                    ))
                                  ) : (
                                    <li className="text-sm text-muted-foreground">No impact specified yet.</li>
                                  )}
                                </ul>
                              </div>
                            )}
                          </div>

                          <div className="flex justify-end mt-6">
                            {fip.showDetailedImpact && (
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => {
                                  setSelectedFIP(fip)
                                  setActiveFIPId(fip.id)
                                }}
                              >
                                View Detailed Impact
                              </Button>
                            )}
                          </div>
                        </div>
                      )}
                    </motion.div>
                  ))
                ) : (
                  <div className="text-center py-12 text-muted-foreground">
                    <p>No FIP specifications available for this upgrade.</p>
                  </div>
                )}
              </div>
            </section>
          </div>
        </div>
      </main>

      {selectedFIP && <FIPImpactModal fip={selectedFIP} isOpen={!!selectedFIP} onClose={() => setSelectedFIP(null)} />}
    </div>
  )
}
