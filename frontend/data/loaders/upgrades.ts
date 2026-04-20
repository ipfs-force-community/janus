import goldenWeekUpgrade from "@/data/upgrades/golden-week-upgrade.json"
import teep from "@/data/upgrades/teep.json"
import tukTuk from "@/data/upgrades/tuk-tuk.json"
import fireHorseUpgrade from "@/data/upgrades/fire-horse-upgrade.json"

export const upgradeLoaders = {
  goldenWeekUpgrade,
  teep,
  "tuk-tuk": tukTuk,
  "fire-horse-upgrade": fireHorseUpgrade,
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
