import fip0077 from "@/data/fips/fip-0077.json"
import frc0102 from "@/data/fips/frc-0102.json"
import fip0103 from "@/data/fips/fip-0103.json"
import fip0106 from "@/data/fips/fip-0106.json"
import fip0101 from "@/data/fips/fip-0101.json"
import frc0108 from "@/data/fips/frc-0108.json"
import fip0109 from "@/data/fips/fip-0109.json"
import fip0097 from "@/data/fips/fip-0097.json"
import fip0081 from "@/data/fips/fip-0081.json"
import fip0085 from "@/data/fips/fip-0085.json"

export const fipLoaders = {
  "fip-0077": fip0077,
  "frc-0102": frc0102,
  "fip-0103": fip0103,
  "fip-0106": fip0106,
  "fip-0101": fip0101,
  "frc-0108": frc0108,
  "fip-0109": fip0109,
  "fip-0097": fip0097,
  "fip-0081": fip0081,
  "fip-0085": fip0085,
}

export const getFIP = (id: string) => {
  return fipLoaders[id.toLowerCase() as keyof typeof fipLoaders] || null
}
