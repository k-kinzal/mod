package gitops

import (
	"github.com/variantdev/mod/pkg/cmdsite"
	"os"
	"strings"
)

type Client struct {
	cmdr    cmdsite.RunCommand
	sh      *cmdsite.CommandSite
	wd      string
	gitPath string
}

func WD(wd string) Option {
	return func(c *Client) {
		c.wd = wd
	}
}
func Commander(cmdr cmdsite.RunCommand) Option {
	return func(c *Client) {
		c.cmdr = cmdr
	}
}

type Option func(*Client)

func New(opt ...Option) *Client {
	c := &Client{}

	for _, o := range opt {
		o(c)
	}

	c.sh = cmdsite.New(cmdsite.RunCmd(c.cmdr))
	c.gitPath = "git"

	return c
}

func (c *Client) Checkout(branch string, has bool) error {
	if has {
		return c.git("checkout", []string{branch})
	}
	return c.git("checkout", []string{"-b", branch})
}

func (c *Client) Add(files ...string) error {
	return c.git("add", files)
}

func (c *Client) Commit(msg string) error {
	return c.git("commit", []string{"-m", msg})
}

func (c *Client) Clone(repo string) error {
	return c.git("clone", []string{repo})
}

func (c *Client) GetCurrentBranch() (string, error) {
	stdout, _, err := c.sh.CaptureStrings(c.gitPath, []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, "\n"), nil
}

func (c *Client) HasBranch(branch string) (bool, error) {
	stdout, _, err := c.sh.CaptureStrings(c.gitPath, []string{"branch", "--list"})
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(stdout, "\n") {
		if strings.TrimSpace(line) == branch {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) GetPushURL(name string) (string, error) {
	stdout, _, err := c.sh.CaptureStrings(c.gitPath, []string{"remote", "get-url", "--push", name})
	if err != nil {
		return "", err
	}
	return stdout, nil
}

func (c *Client) Push(branch string) error {
	return c.git("push", []string{"origin", branch})
}

func (c *Client) DiffExists() bool {
	_, _, err := c.sh.CaptureStrings(c.gitPath, []string{"diff", "--cached", "--exit-code"})
	return err != nil
}

func (c *Client) git(cmd string, args []string) error {
	return c.sh.RunCommand(c.gitPath, append([]string{cmd}, args...), os.Stdout, os.Stderr)
}
