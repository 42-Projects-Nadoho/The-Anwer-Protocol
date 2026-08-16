# World Design

This document provides a comprehensive overview of the game world, including room connections, item distribution, NPCs, and available quests.

## Map Overview

The map consists of an 8-room loop, allowing continuous exploration, with one optional branch leading to the Secret Cave. 

```mermaid
graph TD
    DI["Destiny Islands"] <-->|North / South| TT["Traverse Town"]
    DI <-->|East / West| HB["Hollow Bastion"]
    DI <-->|South / North| SC["Secret Cave (Branch)"]
    
    TT <-->|East / West| W["Wonderland"]
    W <-->|South / North| OC["Olympus Coliseum"]
    OC <-->|West / East| A["Agrabah"]
    A <-->|South / North| HT["Halloween Town"]
    HT <-->|East / West| N["Neverland"]
    N <-->|North / South| HB
```

## Rooms & Exits

* **Destiny Islands**
  * North: Traverse Town
  * East: Hollow Bastion
  * South: Secret Cave
  * *Items:* Potion
  * *NPCs:* Master Yen Sid

* **Secret Cave** (Optional Branch)
  * North: Destiny Islands

* **Traverse Town**
  * South: Destiny Islands
  * East: Wonderland
  * *Items:* Kingdom Key
  * *NPCs:* Leon

* **Wonderland**
  * West: Traverse Town
  * South: Olympus Coliseum

* **Olympus Coliseum**
  * North: Wonderland
  * West: Agrabah
  * *NPCs:* Large Body, Sephiroth

* **Agrabah**
  * East: Olympus Coliseum
  * South: Halloween Town
  * *Items:* Ether

* **Halloween Town**
  * North: Agrabah
  * East: Neverland

* **Neverland**
  * West: Halloween Town
  * North: Hollow Bastion
  * *Items:* Wayfinder

* **Hollow Bastion**
  * South: Neverland
  * West: Destiny Islands
  * *NPCs:* Shadow

---

## Items

| Item ID | Name | Description | Location |
| :--- | :--- | :--- | :--- |
| `potion` | Potion | A restorative item that heals a small amount of HP. | Destiny Islands |
| `ether` | Ether | An item that restores magic power. | Agrabah |
| `keyblade` | Kingdom Key | A mysterious weapon shaped like a giant key. | Traverse Town |
| `wayfinder` | Wayfinder | A star-shaped charm made of thalassa shells. | Neverland |

---

## NPCs

| NPC ID | Name | Role | Location | Description | Stats |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `yen_sid` | Master Yen Sid | Quest Giver | Destiny Islands | A wise and powerful sorcerer, retired Keyblade Master. | - |
| `leon` | Leon | Dialogue | Traverse Town | A stoic warrior wielding a gunblade. | - |
| `shadow_heartless` | Shadow | Enemy | Hollow Bastion | A pureblood Heartless born from the darkness in a heart. | HP: 30, DMG: 10 |
| `large_body` | Large Body | Enemy | Olympus Coliseum | A massive Heartless that blocks attacks with its big belly. | HP: 100, DMG: 25 |
| `sephiroth` | Sephiroth | Enemy | Olympus Coliseum | A legendary SOLDIER with a massive nodachi. | HP: 9999, DMG: 9999 |

---

## Quests

| Quest ID | Name | Type | Target | Reward | Given By |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `find_wayfinder` | Bonds of Friendship | Fetch | `wayfinder` | You feel your heart grow stronger. | Master Yen Sid |
| `defeat_shadow` | Push Back the Darkness | Defeat | `shadow_heartless` | You gained valuable combat experience. | Master Yen Sid |
