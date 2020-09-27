default:

release: toolbin
	deploy/release.bash

toolbin:
	deploy/toolbin.bash

clean:
	rm -rf toolbin frontend power-monitor
