"use client"

import { ThemeToggle } from "@/components/theme-toggle"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { getAllUpgrades } from "@/data/loaders/upgrades"
import upgradesData from "@/data/upgrades.json"
import { motion } from "framer-motion"
import { Calendar, ChevronRight, Clock, GitBranch } from "lucide-react"
import Image from 'next/image'
import { useRouter } from "next/navigation"
import type React from "react"
import { useEffect, useState } from "react"

const { upgradeIds, metadata } = upgradesData
const { statusColors, chainColors, fallbacks } = metadata

type Upgrade = {
  id: string
  name: string
  networkVersion: string
  chain: string
  epochTarget: number
  timeTarget: string
  status: string
  lotusReleaseTag: string
  lotusReleaseUrl: string
  venusReleaseTag: string
  venusReleaseUrl: string
  specs: string[]
  notes: string
  fipIds: string[]
}

function CountdownTimer({ targetTime }: { targetTime: string }) {
  const [timeLeft, setTimeLeft] = useState("")

  useEffect(() => {
    const updateCountdown = () => {
      const now = new Date().getTime()
      const target = new Date(targetTime).getTime()
      const difference = target - now

      if (difference > 0) {
        const days = Math.floor(difference / (1000 * 60 * 60 * 24))
        const hours = Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
        const minutes = Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60))
        const seconds = Math.floor((difference % (1000 * 60)) / 1000)

        setTimeLeft(`${days}d ${hours}h ${minutes}m ${seconds}s`)
      } else {
        setTimeLeft("Completed")
      }
    }

    updateCountdown()
    const interval = setInterval(updateCountdown, 1000)

    return () => clearInterval(interval)
  }, [targetTime])

  return (
    <div className="flex items-center gap-2 text-sm text-muted-foreground">
      <Clock className="h-4 w-4" />
      <span className="font-mono">{timeLeft}</span>
    </div>
  )
}

function UpgradeCard({ upgrade, upgradeId }: { upgrade: Upgrade; upgradeId: string }) {
  const [expanded, setExpanded] = useState(false)
  const router = useRouter()

  const handleCardClick = (e: React.MouseEvent) => {
    if ((e.target as HTMLElement).closest("a, button")) {
      return
    }
    router.push(`/upgrade/${upgradeId}`)
  }

  const safeUpgrade = {
    ...upgrade,
    status: upgrade.status || fallbacks.defaultStatus,
    chain: upgrade.chain || fallbacks.defaultChain,
    notes: upgrade.notes || fallbacks.defaultNotes,
    lotusReleaseTag: upgrade.lotusReleaseTag || fallbacks.defaultReleaseTag,
    lotusReleaseUrl: upgrade.lotusReleaseUrl || fallbacks.defaultReleaseUrl,
    venusReleaseTag: upgrade.venusReleaseTag || fallbacks.defaultReleaseTag,
    venusReleaseUrl: upgrade.venusReleaseUrl || fallbacks.defaultReleaseUrl,
    specs: upgrade.specs || fallbacks.emptySpecs,
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      whileHover={{ y: -2 }}
      className="h-full"
    >
      <Card
        className="cursor-pointer transition-all duration-200 hover:shadow-lg hover:scale-105 h-full flex flex-col"
        onClick={handleCardClick}
      >
        <CardHeader>
          <div className="flex items-start justify-between">
            <div>
              <CardTitle className="text-xl font-bold flex items-center gap-2">
                {safeUpgrade.name}
                <ChevronRight className="h-4 w-4 text-muted-foreground" />
              </CardTitle>
              <p className="text-sm text-muted-foreground mt-1">Network Version {safeUpgrade.networkVersion}</p>
            </div>
            <div className="flex gap-2">
              <Badge className={chainColors[safeUpgrade.chain as keyof typeof chainColors] || chainColors.Mainnet}>
                {safeUpgrade.chain}
              </Badge>
              <Badge className={statusColors[safeUpgrade.status as keyof typeof statusColors] || statusColors.Upcoming}>
                {safeUpgrade.status}
              </Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent className="flex-1 flex flex-col">
          <div className="space-y-4 flex-1">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Target Epoch</p>
                <p className="font-mono text-lg">{safeUpgrade.epochTarget?.toLocaleString() || "TBD"}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">{
                  safeUpgrade.status === "Finalized" ? "Finalized Time" : "Estimated Time"
                }</p>
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4" />
                  <span className="text-sm">
                    {safeUpgrade.timeTarget
                      ? new Date(safeUpgrade.timeTarget).toLocaleDateString("en-US", {
                        year: "numeric",
                        month: "short",
                        day: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })
                      : "TBD"}
                  </span>
                </div>
              </div>
            </div>

            {safeUpgrade.status === "Upcoming" && safeUpgrade.timeTarget && (
              <CountdownTimer targetTime={safeUpgrade.timeTarget} />
            )}
            <div>
              <p className="text-sm font-medium text-muted-foreground mb-2">Release Tags</p>
              <div className="flex items-center gap-6">
                <div className="flex items-center gap-2">
                  <GitBranch className="h-4 w-4" />
                  <span className="bg-muted px-2 py-1 rounded text-sm">
                    {safeUpgrade.lotusReleaseTag}
                  </span>
                </div>

                <div className="flex items-center gap-2">
                  <GitBranch className="h-4 w-4" />
                  <span className="bg-muted px-2 py-1 rounded text-sm">
                    {safeUpgrade.venusReleaseTag}
                  </span>
                </div>
              </div>
            </div>

            <div className="flex-1">
              <p className="text-sm font-medium text-muted-foreground mb-2">FIP Specifications</p>
              <div className="flex flex-wrap gap-2">
                {safeUpgrade.specs.length > 0 ? (
                  safeUpgrade.specs.map((spec) => (
                    <Badge key={spec} variant="outline" className="text-xs">
                      {spec}
                    </Badge>
                  ))
                ) : (
                  <span className="text-sm text-muted-foreground">No FIP specifications available</span>
                )}
              </div>
            </div>
          </div>

          <div className="text-xs text-muted-foreground border-t pt-3 mt-4 text-center hover:text-primary">
            Click to view detailed FIP analysis →
          </div>
        </CardContent>
      </Card>
    </motion.div>
  )
}

export default function FilecoinUpgrades() {
  const [statusFilter, setStatusFilter] = useState("all")
  const [searchQuery, setSearchQuery] = useState("")
  const [upgrades, setUpgrades] = useState<{ upgrade: Upgrade; upgradeId: string }[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const loadAllUpgrades = () => {
      setLoading(true)
      const loadedUpgrades = getAllUpgrades()
      setUpgrades(loadedUpgrades)
      setLoading(false)
    }

    loadAllUpgrades()
  }, [])

  const filteredUpgrades = upgrades.filter(({ upgrade }) => {
    const matchesStatus = statusFilter === "all" || upgrade.status === statusFilter
    const matchesSearch =
      searchQuery === "" ||
      upgrade.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      upgrade.networkVersion.toLowerCase().includes(searchQuery.toLowerCase())

    return matchesStatus && matchesSearch
  })

  const statuses = Array.from(new Set(upgrades.map(({ upgrade }) => upgrade.status).filter(Boolean)))

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading upgrades...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="text-center flex-1">
              <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
                Filecoin Upgrades
              </h1>
              <p className="text-muted-foreground mt-2">
                Track Filecoin Mainnet upgrades and their progress
              </p>
            </div>
            <ThemeToggle />
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row gap-4 mb-8">
          <Input
            placeholder="Search by name or version (e.g., Teep, nv25)..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="md:flex-1"
          />
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="md:w-48">
              <SelectValue placeholder="Filter by Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Statuses</SelectItem>
              {statuses.map((status) => (
                <SelectItem key={status} value={status}>
                  {status}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="mt-8">
          {filteredUpgrades.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-muted-foreground">
                No upgrades found matching your criteria.
              </p>
            </div>
          ) : (
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {filteredUpgrades.map(({ upgrade, upgradeId }) => (
                <UpgradeCard key={upgradeId} upgrade={upgrade} upgradeId={upgradeId} />
              ))}
            </div>
          )}
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t bg-card mt-auto">
        <div className="container mx-auto px-4 py-6 flex flex-col items-center gap-2 text-center">
          <p className="text-sm text-muted-foreground">Powered by Venus</p>
          <div className="flex items-center gap-6">
            <a
              href="https://filecoinproject.slack.com/archives/CEHHJNJS3"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center text-sm text-muted-foreground hover:text-primary hover:underline"
            >
              <Image
                src="/icons/slack.png"
                alt="Slack"
                width={16}
                height={16}
                className="mr-2"
              />
              Slack
            </a>
            <a
              href="https://github.com/filecoin-project/venus"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center text-sm text-muted-foreground hover:text-primary hover:underline"
            >
              <Image
                src="/icons/github.png"
                alt="GitHub"
                width={16}
                height={16}
                className="mr-2"
              />
              GitHub
            </a>
          </div>
        </div>
      </footer>
    </div>
  )
}
