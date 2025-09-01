import mercury from "@/data/upgrades/mercury.json"
import teep from "@/data/upgrades/teep.json"
import tukTuk from "@/data/upgrades/tuk-tuk.json"

export const upgradeLoaders = {
  mercury,
  teep,
  "tuk-tuk": tukTuk,
}

export const getUpgrade = (id: string) => {
  return upgradeLoaders[id as keyof typeof upgradeLoaders] || null
}

export const getAllUpgrades = () => {
  return Object.entries(upgradeLoaders).map(([id, upgrade]) => ({
    upgrade,
    upgradeId: id,
  }))
}
