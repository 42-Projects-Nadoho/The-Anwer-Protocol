# Group Contributions

The Answer Protocol (TAP) is a collaborative project created by a team of two developers. To maximize efficiency and ensure a robust implementation, the responsibilities were clearly divided between the group members.

## Team Members and Responsibilities

| Memeber | Role |
| :--- | :--- |
| nadoho | Server Implementation and CLI Client |
| spacotto | GUI Client and World Design |

### Server Implementation and CLI Client

| Task | Description |
| :--- | :--- |
| **Server Architecture** | Designed and implemented the core TCP server capable of handling multiple concurrent connections using Go's goroutines. |
| **Protocol Management** | Implemented the server-side parsing and routing of the TAP protocol commands according to the RFC specifications. |
| **State Synchronization** | Handled the synchronization of the shared world state and managed real-time broadcasts to players. |
| **CLI Client** | Built the lightweight, text-based command-line interface, ensuring it properly formats commands, parses server events asynchronously, and provides a smooth experience for terminal users. |

### GUI Client and World Design
**Responsibilities:** GUI Client and World Design

| Task | Description |
| :--- | :--- |
| **GUI Client** | Designed and developed the rich graphical user interface. Ensured that the interface remains responsive while receiving real-time events and provides dedicated sections for room details, chat scopes, and inventory management. |
| **World Design** | Architected the layout of the virtual world, ensuring all requirements (loops, optional branches) were met. |
| **Data Implementation** | Authored the `world.yaml` configuration file, populating the world with interconnected rooms, engaging NPCs (quest-givers, enemies, dialogue characters), and discoverable items. |
| **Quest & Combat Integration** | Tailored the placement of specific NPCs and items to support the flow of the combat and quest progression mechanics. |
